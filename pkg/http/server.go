// Package http provides utilities for handling HTTP requests and responses.
// It includes types and functions for creating custom handlers, applying middleware,
// and managing HTTP endpoints.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shandysiswandi/gostarter/pkg/errs"
)

// ServeHandlerFunc defines a function that processes an HTTP request within a given context
// and returns a response of any type along with an error, if one occurred.
// The context can be used to carry request-scoped values, cancellation signals, and deadlines.
type ServeHandlerFunc func(context.Context, *http.Request) (any, error)

// Middleware defines a function type that wraps an http.Handler, adding additional
// functionality such as logging, authentication, or request modification.
// Middleware is applied in the order they are added.
type Middleware func(http.Handler) http.Handler

// StatusCoder is an interface for custom responses that allows them to specify their
// HTTP status code. Types implementing this interface can control the HTTP status code
// returned to the client.
type StatusCoder interface {
	// StatusCode returns the HTTP status code for the response.
	StatusCode() int
}

// ServeOption represents a configuration option that can be applied to a Serve instance.
// It allows customization of the Serve instance through the Apply method.
type ServeOption interface {
	Apply(s *Serve)
}

// ServeOptionFunc is a function type that implements the ServeOption interface.
// It allows functions to be used as configuration options when creating a Serve instance.
type ServeOptionFunc func(*Serve)

// Apply calls the ServeOptionFunc, applying the configuration to the Serve instance.
func (o ServeOptionFunc) Apply(s *Serve) {
	o(s)
}

// Serve provides methods for creating and managing HTTP endpoints.
// It supports middleware chaining and customizable request handling.
type Serve struct {
	mws []Middleware
}

// NewServe creates and returns a new Serve instance with the provided options.
// Serve instances manage the global middleware stack and create individual endpoints.
func NewServe(options ...ServeOption) *Serve {
	s := &Serve{}

	for _, o := range options {
		o.Apply(s)
	}

	return s
}

// Endpoint creates and returns a new Endpoint with the specified handler and optional middleware.
// The endpoint will execute the middleware stack defined in Serve, followed by the provided middleware
// for the specific endpoint.
func (s *Serve) Endpoint(h ServeHandlerFunc, mws ...Middleware) *Endpoint {
	return &Endpoint{
		h:   h,
		mws: append(s.mws, mws...),
	}
}

// WithMiddlewares is a helper function that creates a ServeOption to add middleware to the Serve instance.
// It appends the provided middleware to the Serve's global middleware stack.
func WithMiddlewares(mws ...Middleware) ServeOption {
	return ServeOptionFunc(func(s *Serve) {
		s.mws = append(s.mws, mws...)
	})
}

// Endpoint represents an HTTP endpoint with an associated handler and middleware stack.
// It is responsible for processing HTTP requests, executing middleware, and encoding responses.
type Endpoint struct {
	h   ServeHandlerFunc
	mws []Middleware
}

// baseHandler is an internal method that processes incoming requests by executing the handler
// and encoding the response. If an error occurs, it encodes the error response.
func (e *Endpoint) baseHandler(w http.ResponseWriter, r *http.Request) {
	res, err := e.h(r.Context(), r)
	if err != nil {
		e.errorEncoder(w, err)

		return
	}

	if err := e.responseEncoder(w, res); err != nil {
		e.errorEncoder(w, err)

		return
	}
}

// chainMiddleware chains the provided middleware functions around an http.Handler.
// The middleware are applied in the order they are provided, with the first middleware
// wrapping the handler and each subsequent middleware wrapping the previous one.
func (e *Endpoint) chainMiddleware(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

// ServeHTTP handles incoming HTTP requests by applying the middleware chain and
// then invoking the endpoint's handler. It is the entry point for the endpoint's
// request processing logic.
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := e.chainMiddleware(http.HandlerFunc(e.baseHandler), e.mws...)
	handler.ServeHTTP(w, r)
}

// errorEncoder encodes the error response as a JSON object with an "error" field
// containing the error message. It also sets the appropriate HTTP status code based
// on the error, defaulting to 500 Internal Server Error if no custom code is provided.
func (e *Endpoint) errorEncoder(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.errCode(err))

	//nolint:errcheck,errchkjson // ignoring error as we're encoding JSON
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

// responseEncoder encodes the response as a JSON object. It sets the HTTP status code
// based on the response's StatusCoder implementation, or defaults to 200 OK.
// If the status code is 204 No Content, the response body is omitted.
func (e *Endpoint) responseEncoder(w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusOK
	if sc, ok := response.(StatusCoder); ok {
		code = sc.StatusCode()
	}

	if code == http.StatusNoContent {
		w.WriteHeader(code)

		return nil
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	if err := encoder.Encode(response); err != nil {
		return err
	}

	w.WriteHeader(code)

	_, _ = buf.WriteTo(w) //nolint:errcheck // it will not error, maybe

	return nil
}

// errCode determines the appropriate HTTP status code based on the provided error.
// It checks if the error is of type *errs.Error and maps it to a specific HTTP status code.
// Defaults to 500 Internal Server Error if the error type is not recognized.
func (e *Endpoint) errCode(err error) int {
	var ee *errs.Error
	if !errors.As(err, &ee) {
		return http.StatusInternalServerError
	}

	if ok := errs.IsServerError(ee); ok {
		return http.StatusInternalServerError
	}

	if ok := errs.IsValidationError(ee); ok {
		return http.StatusBadRequest
	}

	switch ee.Code() {
	case errs.CodeInvalidInput:
		return http.StatusBadRequest
	case errs.CodeNotFound:
		return http.StatusNotFound
	case errs.CodeConflict:
		return http.StatusConflict
	case errs.CodeUnauthorized:
		return http.StatusUnauthorized
	case errs.CodeForbidden:
		return http.StatusForbidden
	case errs.CodeTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
