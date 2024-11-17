package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	msgInternal  = "internal server error"
	valAppJSONCT = "application/json; charset=utf-8"
)

type Vendor int

const (
	VendorDefault Vendor = iota
	VendorHTTPRouter
)

// Handler defines the type for endpoint handlers with context, request, and response writer.
type Handler func(context.Context, *http.Request) (any, error)

// Middleware defines the type for middleware functions.
type Middleware func(http.Handler) http.Handler

// StatusCoder is an interface for custom responses that allows them to specify their
// HTTP status code. Types implementing this interface can control the HTTP status code
// returned to the client.
type StatusCoder interface {
	StatusCode() int
}

type Router struct {
	router           *httprouter.Router
	vendor           Vendor
	notFound         http.Handler
	methodNotAllowed http.Handler
	resultCodec      func(context.Context, http.ResponseWriter, any) error
	errorCodec       func(context.Context, http.ResponseWriter, error)
}

func New() *Router {
	return &Router{
		router:           httprouter.New(),
		vendor:           VendorHTTPRouter,
		notFound:         http.HandlerFunc(defaultNotFound),
		methodNotAllowed: http.HandlerFunc(defaultMethodNotAllowedh),
		resultCodec:      defaultResultCodec,
		errorCodec:       defaultErrorCodec,
	}
}

func (r *Router) Endpoint(method, path string, h Handler, mws ...Middleware) {
	hm := chainMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		res, err := h(ctx, req)
		if err != nil {
			r.errorCodec(ctx, w, err)

			return
		}

		if err := r.resultCodec(ctx, w, res); err != nil {
			http.Error(w, msgInternal, http.StatusInternalServerError)

			return
		}
	}), mws...)

	r.router.Handler(method, path, hm)
}

func (r *Router) Native(method, path string, h http.Handler, mws ...Middleware) {
	r.router.Handler(method, path, chainMiddleware(h, mws...))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.NotFound = r.notFound
	r.router.MethodNotAllowed = r.methodNotAllowed
	r.router.ServeHTTP(w, req)
}

func chainMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

// defaultResultCodec is the default implementation for encoding successful responses.
// It sets the content type to JSON and writes a status code of 200 OK.
// The response body contains the encoded `data` in a JSON object.
func defaultResultCodec(_ context.Context, w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", valAppJSONCT)

	code := http.StatusOK
	if sc, ok := data.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)

		return nil
	}

	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(data)
}

// defaultErrorCodec is the default implementation for encoding error responses.
// It sets the content type to JSON and writes a status code of 500 Internal Server Error.
// The response body contains the error message in a JSON object.
func defaultErrorCodec(_ context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", valAppJSONCT)

	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)

	err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if err != nil {
		return
	}
}

func defaultNotFound(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", valAppJSONCT)
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(map[string]string{"error": "endpoint not found"})
	if err != nil {
		http.Error(w, msgInternal, http.StatusInternalServerError)
	}
}

func defaultMethodNotAllowedh(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", valAppJSONCT)
	w.WriteHeader(http.StatusMethodNotAllowed)
	err := json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	if err != nil {
		http.Error(w, msgInternal, http.StatusInternalServerError)
	}
}
