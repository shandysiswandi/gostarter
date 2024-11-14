// Package serve provides utilities for handling HTTP requests and responses.
// It includes types and functions for creating custom handlers, applying middleware,
// and managing HTTP endpoints.
package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// StatusCoder is an interface for custom responses that allows them to specify their
// HTTP status code. Types implementing this interface can control the HTTP status code
// returned to the client.
type StatusCoder interface {
	// StatusCode returns the HTTP status code for the response.
	StatusCode() int
}

// Serve provides methods for creating and managing HTTP endpoints.
// It supports middleware chaining and customizable request handling.
type Serve struct {
	mws         []func(http.Handler) http.Handler
	resultCodec func(context.Context, http.ResponseWriter, any) error
	errorCodec  func(context.Context, http.ResponseWriter, error)
}

// WithMiddlewares is a helper function to add middleware to the Serve instance.
// It appends the provided middleware to the Serve's global middleware stack.
func WithMiddlewares(mws ...func(http.Handler) http.Handler) func(*Serve) {
	return func(s *Serve) {
		s.mws = append(s.mws, mws...)
	}
}

// WithResultCodec sets the result codec function for encoding responses.
// The result codec is used to encode responses of type `any` before sending them to the client.
func WithResultCodec(rc func(context.Context, http.ResponseWriter, any) error) func(*Serve) {
	return func(s *Serve) {
		s.resultCodec = rc
	}
}

// WithErrorCodec sets the error codec function for encoding errors.
// The error codec is used to encode errors before sending them to the client.
func WithErrorCodec(ec func(context.Context, http.ResponseWriter, error)) func(*Serve) {
	return func(s *Serve) {
		s.errorCodec = ec
	}
}

// defaultResultCodec is the default implementation for encoding successful responses.
// It sets the content type to JSON and writes a status code of 200 OK.
// The response body contains the encoded `data` in a JSON object.
func defaultResultCodec(_ context.Context, w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusOK
	if sc, ok := data.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)

		return nil
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	w.WriteHeader(code)

	_, _ = buf.WriteTo(w) //nolint:errcheck // ignore for this, it never error

	return nil
}

// defaultErrorCodec is the default implementation for encoding error responses.
// It sets the content type to JSON and writes a status code of 500 Internal Server Error.
// The response body contains the error message in a JSON object.
func defaultErrorCodec(_ context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)

	//nolint:errcheck,errchkjson // ignore for this, it never error
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// New creates and returns a new Serve instance with the provided options.
// The new Serve instance is initialized with default result and error codecs.
// Serve instances manage the global middleware stack and create individual endpoints.
func New(opts ...func(*Serve)) *Serve {
	s := &Serve{errorCodec: defaultErrorCodec, resultCodec: defaultResultCodec}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Endpoint creates and returns a new Endpoint with the specified handler and optional middleware.
// The endpoint will execute the middleware stack defined in Serve,
// followed by the provided middleware for the specific endpoint.
// The Endpoint will use the result and error codecs defined in Serve for encoding responses and errors.
func (s *Serve) Endpoint(h func(context.Context, *http.Request) (any, error),
	mws ...func(http.Handler) http.Handler,
) *Endpoint {
	return &Endpoint{
		h:           h,
		mws:         append(s.mws, mws...),
		resultCodec: s.resultCodec,
		errorCodec:  s.errorCodec,
	}
}
