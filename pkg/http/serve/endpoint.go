// Package serve provides utilities for handling HTTP requests and responses.
// It includes types and functions for creating custom handlers, applying middleware,
// and managing HTTP endpoints.
package serve

import (
	"context"
	"net/http"
)

// Endpoint represents an HTTP endpoint with an associated handler and middleware stack.
// It is responsible for processing HTTP requests, executing middleware, and encoding responses.
type Endpoint struct {
	h           func(context.Context, *http.Request) (any, error)
	mws         []func(http.Handler) http.Handler
	resultCodec func(context.Context, http.ResponseWriter, any) error
	errorCodec  func(context.Context, http.ResponseWriter, error)
}

// handler is an internal method that processes incoming requests by executing the handler
// and encoding the response. If an error occurs, it encodes the error response.
func (e *Endpoint) handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := e.h(ctx, r)
	if err != nil {
		e.errorCodec(ctx, w, err)

		return
	}

	if err := e.resultCodec(ctx, w, res); err != nil {
		e.errorCodec(ctx, w, err)

		return
	}
}

// chainMiddleware chains the provided middleware functions around an http.Handler.
// The middleware are applied in the order they are provided, with the first middleware
// wrapping the handler and each subsequent middleware wrapping the previous one.
func (e *Endpoint) chainMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}

	return h
}

// ServeHTTP handles incoming HTTP requests by applying the middleware chain and
// then invoking the endpoint's handler. It is the entry point for the endpoint's
// request processing logic.
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.chainMiddleware(http.HandlerFunc(e.handler), e.mws...).ServeHTTP(w, r)
}
