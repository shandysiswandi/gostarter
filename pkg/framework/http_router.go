package framework

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	valAppJSONCT   = "application/json; charset=utf-8"
	keyContentType = "Content-Type"
)

// Handler defines the type for endpoint handlers with context, request, and response writer.
type Handler func(Context) (any, error)

// StatusCoder is an interface for custom responses that allows them to specify their
// HTTP status code. Types implementing this interface can control the HTTP status code
// returned to the client.
type StatusCoder interface {
	StatusCode() int
}

// Router manages HTTP routes, handlers, and middleware.
type Router struct {
	router           *httprouter.Router
	notFound         http.Handler
	methodNotAllowed http.Handler
	resultCodec      func(context.Context, http.ResponseWriter, any)
	errorCodec       func(context.Context, http.ResponseWriter, error)
}

// New creates and initializes a new Router instance with default configurations.
func New() *Router {
	return &Router{
		router:           httprouter.New(),
		notFound:         http.HandlerFunc(defaultNotFound),
		methodNotAllowed: http.HandlerFunc(defaultMethodNotAllowed),
		resultCodec:      defaultResultCodec,
		errorCodec:       defaultErrorCodec,
	}
}

// Endpoint registers a new route with a specific HTTP method, path, handler, and optional middleware.
func (r *Router) Endpoint(method, path string, h Handler, mws ...Middleware) {
	wh := r.wrapHandler(h)
	cm := Chain(wh, mws...)
	r.router.Handler(method, path, cm)
}

// Native registers a native HTTP handler with a specific method, path, and optional middleware.
func (r *Router) Native(method, path string, h http.Handler, mws ...Middleware) {
	r.router.Handler(method, path, Chain(h, mws...))
}

// ServeHTTP dispatches the request to the appropriate route handler or middleware.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.NotFound = r.notFound
	r.router.MethodNotAllowed = r.methodNotAllowed
	r.router.ServeHTTP(w, req)
}

// wrapHandler wraps a Handler to integrate it with the Router's context and error/result codecs.
func (r *Router) wrapHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cc := &RouterCtx{r: req}

		res, err := h(cc)
		if err != nil {
			r.errorCodec(req.Context(), w, err)

			return
		}

		r.resultCodec(req.Context(), w, res)
	})
}

// defaultResultCodec encodes successful responses into JSON format.
// It sets the content type to JSON and writes the appropriate HTTP status code.
func defaultResultCodec(_ context.Context, w http.ResponseWriter, data any) {
	w.Header().Set(keyContentType, valAppJSONCT)

	code := http.StatusOK
	if sc, ok := data.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	w.WriteHeader(code)
	writeJSON(w, data)
}

// defaultErrorCodec encodes error responses into JSON format.
// It sets the content type to JSON and writes the appropriate HTTP status code.
func defaultErrorCodec(_ context.Context, w http.ResponseWriter, err error) {
	w.Header().Set(keyContentType, valAppJSONCT)

	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)
	writeJSON(w, map[string]string{"error": err.Error()})
}

// defaultNotFound is the default handler for unregistered routes.
// It responds with a 404 Not Found status and a JSON error message.
func defaultNotFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(keyContentType, valAppJSONCT)
	w.WriteHeader(http.StatusNotFound)
	writeJSON(w, map[string]string{"error": "endpoint not found"})
}

// defaultMethodNotAllowed is the default handler for unsupported HTTP methods.
// It responds with a 405 Method Not Allowed status and a JSON error message.
func defaultMethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set(keyContentType, valAppJSONCT)
	w.WriteHeader(http.StatusMethodNotAllowed)
	writeJSON(w, map[string]string{"error": "method not allowed"})
}

// writeJSON is a utility function to write a JSON response to an http.ResponseWriter.
// It encodes the provided data as JSON and writes it to the response body.
// If an error occurs during encoding, it responds with an internal server error (500).
func writeJSON(w http.ResponseWriter, data any) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println("json.NewEncoder(w).Encode(data)", err)
	}
}
