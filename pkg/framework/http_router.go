package framework

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler defines the type for endpoint handlers with context, request, and response writer.
type Handler func(Context) (any, error)

// StatusCoder is an interface for custom responses that allows them to specify their
// HTTP status code. Types implementing this interface can control the HTTP status code
// returned to the client.
type StatusCoder interface {
	StatusCode() int
}

type Router struct {
	hr          *httprouter.Router
	resultCodec func(context.Context, http.ResponseWriter, any)
	errorCodec  func(context.Context, http.ResponseWriter, error)
}

func NewRouter() *Router {
	return &Router{
		hr: &httprouter.Router{
			HandleMethodNotAllowed: true,
			SaveMatchedRoutePath:   true,
			NotFound:               http.HandlerFunc(defaultNotFound),
			MethodNotAllowed:       http.HandlerFunc(defaultMethodNotAllowed),
		},
		resultCodec: defaultResultCodec,
		errorCodec:  defaultErrorCodec,
	}
}

func (r *Router) Endpoint(method, path string, h Handler, mws ...Middleware) {
	r.hr.Handler(method, path, Chain(http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
		rr.Header.Set("X-Actual-Path", httprouter.ParamsFromContext(rr.Context()).MatchedRoutePath())
		cc := &RouterCtx{r: rr}

		res, err := h(cc)
		if err != nil {
			r.errorCodec(rr.Context(), w, err)

			return
		}

		r.resultCodec(rr.Context(), w, res)
	}), mws...))
}

func (r *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	r.hr.HandlerFunc(method, path, handler)
}

func (r *Router) Handler(method, path string, handler http.Handler) {
	r.hr.Handler(method, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.hr.ServeHTTP(w, req)
}

// defaultResultCodec encodes successful responses into JSON format.
// It sets the content type to JSON and writes the appropriate HTTP status code.
func defaultResultCodec(_ context.Context, w http.ResponseWriter, data any) {
	code := http.StatusOK
	if sc, ok := data.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	writeJSON(w, data, code)
}

// defaultErrorCodec encodes error responses into JSON format.
// It sets the content type to JSON and writes the appropriate HTTP status code.
func defaultErrorCodec(_ context.Context, w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	writeJSON(w, map[string]string{"error": err.Error()}, code)
}

// defaultNotFound is the default handler for unregistered routes.
// It responds with a 404 Not Found status and a JSON error message.
func defaultNotFound(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]string{"error": "endpoint not found"}, http.StatusNotFound)
}

// defaultMethodNotAllowed is the default handler for unsupported HTTP methods.
// It responds with a 405 Method Not Allowed status and a JSON error message.
func defaultMethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
}

// writeJSON is a utility function to write a JSON response to an http.ResponseWriter.
// It encodes the provided data as JSON and writes it to the response body.
// If an error occurs during encoding, it responds with an internal server error (500).
func writeJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println("json.NewEncoder(w).Encode(data)", err)
	}
}
