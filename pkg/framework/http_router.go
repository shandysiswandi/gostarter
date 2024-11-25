package framework

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

// nodeType represents the type of node in the radix tree.
type nodeType uint8

// Constants representing different node types.
const (
	static   nodeType = iota // Default: Static route
	root                     // Root node (for the root of the router)
	param                    // Dynamic parameter (e.g., :id)
	catchAll                 // Catch-all wildcard (e.g., *path)
)

// Param represents a single route parameter.
type Param struct {
	Key   string // The name of the parameter (e.g., "id")
	Value string // The value of the parameter (e.g., "123")
}

// Params is a slice of Param, used for passing route parameters.
type Params []Param

func (ps Params) Key(key string) string {
	for _, p := range ps {
		if p.Key == key {
			return p.Value
		}
	}

	return ""
}

type paramContextKey struct{}

// GetParams retrieves route parameters from the request context.
func GetParams(ctx context.Context) Params {
	p, _ := ctx.Value(paramContextKey{}).(Params)

	return p
}

// node represents a node in the radix tree.
type node struct {
	path     string              // Static part of the path or the parameter name
	nodeType nodeType            // Type of the node: static, param, catchAll, etc.
	children map[string]*node    // Child nodes for static segments
	handler  http.Handler        // Handler function for this route
	param    string              // Named parameter (e.g., :id)
	wildcard *node               // Catch-all route node (*path)
	methods  map[string]struct{} // Store allowed HTTP methods (GET, POST, etc.)
}

func newNode(path string, nodeType nodeType) *node {
	return &node{
		path:     path,
		nodeType: nodeType,
		children: make(map[string]*node),
	}
}

// splitPath splits a path string into segments.
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

func countParams(path string) uint8 {
	var n uint8
	for i := range len(path) - 1 {
		if path[i] == ':' || path[i] == '*' {
			n++
		}
	}

	return n
}

// Router is the main struct for the router, containing the root node and a cache.
type Router struct {
	root                    *node
	cache                   map[string]*node // Cache for optimizing route lookups
	paramsPool              sync.Pool        // Pool for route parameters
	maxParams               uint8
	notFoundHandler         http.HandlerFunc
	methodNotAllowedHandler http.HandlerFunc
	resultCodec             func(context.Context, http.ResponseWriter, any)
	errorCodec              func(context.Context, http.ResponseWriter, error)
}

// NewRouter initializes a new RadixRouter.
func NewRouter() *Router {
	return &Router{
		root:                    newNode("", root), // Root node of the radix tree
		cache:                   make(map[string]*node),
		notFoundHandler:         defaultNotFound,
		methodNotAllowedHandler: defaultMethodNotAllowed,
		resultCodec:             defaultResultCodec,
		errorCodec:              defaultErrorCodec,
	}
}

func (r *Router) Endpoint(method, path string, h Handler, mws ...Middleware) {
	wh := http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
		cc := &RouterCtx{r: rr}

		res, err := h(cc)
		if err != nil {
			r.errorCodec(rr.Context(), w, err)

			return
		}

		r.resultCodec(rr.Context(), w, res)
	})
	cm := Chain(wh, mws...)
	r.addRoute(method, path, cm)
}

// Handler adds a route handler to the router for a specific method and path.
func (r *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	r.addRoute(method, path, handler)
}

func (r *Router) Handler(method, path string, handler http.Handler) {
	r.addRoute(method, path, handler)
}

func (r *Router) addRoute(method, path string, handler http.Handler) {
	segments := splitPath(path)
	current := r.root

	paramsCount := countParams(path)
	if paramsCount > r.maxParams {
		r.maxParams = paramsCount
	}

	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)

			return &ps
		}
	}

	for _, segment := range segments {
		switch {
		case strings.HasPrefix(segment, ":"):
			if current.param == "" {
				current.param = segment[1:] // Store parameter name
				current.children[":"] = newNode(segment, param)
			}
			current = current.children[":"]

		case segment == "*":
			if current.wildcard == nil {
				current.wildcard = newNode(segment, catchAll)
			}
			current = current.wildcard
		default:
			if _, exists := current.children[segment]; !exists {
				current.children[segment] = newNode(segment, static)
			}
			current = current.children[segment]
		}
	}

	// Set the handler and method for the last node
	current.handler = handler
	if current.methods == nil {
		current.methods = make(map[string]struct{})
	}
	current.methods[method] = struct{}{} // Add method to allowed methods
}

// getRoute looks up a route in the radix tree.
func (r *Router) getRoute(path string, method string) (*node, *Params, bool, bool) {
	// Check the cache first
	if cachedNode, found := r.cache[path]; found {
		if _, methodAllowed := cachedNode.methods[method]; methodAllowed {
			return cachedNode, nil, true, true
		}

		return cachedNode, nil, true, false // Method not allowed
	}

	segments := splitPath(path)
	current := r.root
	params := r.paramsPool.Get().(*Params)

OuterLoop:
	for _, segment := range segments {
		var child *node
		var exists bool
		if child, exists = current.children[segment]; exists {
			current = child

			continue
		}

		switch {
		case current.param != "":
			*params = append(*params, Param{Key: current.param, Value: segment})
			current = current.children[":"]
		case current.wildcard != nil:
			*params = append(*params, Param{Key: "*", Value: strings.Join(segments, "/")})
			current = current.wildcard

			break OuterLoop
		default:
			r.paramsPool.Put(params)

			return nil, nil, false, false // Not found
		}
	}

	if _, ok := current.methods[method]; !ok {
		r.paramsPool.Put(params)

		return current, params, true, false // Method not allowed
	}

	// Cache the result for future lookups
	r.cache[path] = current

	return current, params, true, true // Found and method allowed
}

// ServeHTTP is the method to handle HTTP requests.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	node, params, found, methodAllowed := r.getRoute(req.URL.Path, req.Method)

	if found && methodAllowed && params != nil {
		ctx := context.WithValue(req.Context(), paramContextKey{}, *params)
		node.handler.ServeHTTP(w, req.WithContext(ctx))
		r.paramsPool.Put(params)

		return
	}

	if found && methodAllowed && params == nil {
		node.handler.ServeHTTP(w, req)

		return
	}

	if found && !methodAllowed {
		r.methodNotAllowedHandler.ServeHTTP(w, req)

		return
	}

	r.notFoundHandler.ServeHTTP(w, req)
}
