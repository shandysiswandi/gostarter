package framework

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/julienschmidt/httprouter"
)

// Context defines an interface for working with HTTP requests and their context.
// It provides methods to access the request, headers, body, query parameters, and route parameters.
type Context interface {
	// Context returns the context.Context of the current request.
	Context() context.Context

	// Request returns the underlying *http.Request.
	Request() *http.Request

	// Header returns the HTTP headers associated with the request.
	Header() http.Header

	// Body returns the body of the HTTP request as an io.ReadCloser.
	Body() io.ReadCloser

	// Query retrieves the first value associated with the given query key.
	Query(key string) string

	// Queries retrieves all values associated with the given query key.
	Queries(key string) []string

	// Param retrieves the value of a route parameter by its key.
	Param(key string) string
}

// RouterCtx is an implementation of the Context interface.
// It wraps an *http.Request and provides utility methods to access
// its context, headers, body, query parameters, and route parameters.
type RouterCtx struct {
	r *http.Request
}

// Context returns the context.Context of the underlying HTTP request.
func (rc *RouterCtx) Context() context.Context {
	return rc.r.Context()
}

// Request returns the underlying *http.Request.
func (rc *RouterCtx) Request() *http.Request {
	return rc.r
}

// Header returns the HTTP headers of the request.
func (rc *RouterCtx) Header() http.Header {
	return rc.r.Header
}

// Body returns the body of the HTTP request as an io.ReadCloser.
func (rc *RouterCtx) Body() io.ReadCloser {
	return rc.r.Body
}

// Query retrieves the first value associated with the given query key.
func (rc *RouterCtx) Query(key string) string {
	return rc.r.URL.Query().Get(key)
}

// Queries retrieves all values associated with the given query key.
func (rc *RouterCtx) Queries(key string) []string {
	return rc.r.URL.Query()[key]
}

// Param retrieves the value of a route parameter by its key.
func (rc *RouterCtx) Param(key string) string {
	return httprouter.ParamsFromContext(rc.Request().Context()).ByName(key)
}

// TestCtx is a helper type for testing HTTP requests and contexts.
// It allows setting custom route parameters, query parameters, and headers.
type TestCtx struct {
	r       *http.Request
	param   map[string]string
	queries map[string][]string
	headers map[string][]string
	mu      sync.RWMutex
}

// NewTestContext creates a new TestDefault instance with a given HTTP method,
// target URL, and optional body.
func NewTestContext(method string, target string, body io.Reader) *TestCtx {
	r := httptest.NewRequest(method, target, body)

	return &TestCtx{
		r:       r,
		param:   make(map[string]string),
		queries: make(map[string][]string),
		headers: make(map[string][]string),
	}
}

// SetParam adds or updates a route parameter key-value pair.
func (tc *TestCtx) SetParam(key, value string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.param[key] = value
}

// SetQuery adds or updates a query parameter with the given key and values.
func (tc *TestCtx) SetQuery(key string, values ...string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.queries[key] = append(tc.queries[key], values...)
}

// SetHeader adds or updates a header with the given key and values.
func (tc *TestCtx) SetHeader(key string, values ...string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.headers[key] = append(tc.headers[key], values...)
}

// Build constructs a RouterCtx instance from the configured TestDefault.
// It sets the query parameters, headers, and route parameters on the HTTP request.
func (tc *TestCtx) Build() *RouterCtx {
	// Set query parameters
	q := tc.r.URL.Query()
	for key, query := range tc.queries {
		for _, value := range query {
			q.Add(key, value)
		}
	}
	tc.r.URL.RawQuery = q.Encode()

	// Set headers
	for key, header := range tc.headers {
		for _, value := range header {
			tc.r.Header.Add(key, value)
		}
	}

	// Set route parameters
	var params httprouter.Params
	for key, value := range tc.param {
		params = append(params, httprouter.Param{
			Key:   key,
			Value: value,
		})
	}
	ctx := context.WithValue(tc.r.Context(), httprouter.ParamsKey, params)

	return &RouterCtx{r: tc.r.WithContext(ctx)}
}
