package router

// import (
// 	"context"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/gorilla/mux"
// 	"github.com/julienschmidt/httprouter"
// )

// // https://github.com/julienschmidt/httprouter
// // https://github.com/go-chi/chi
// // https://github.com/gorilla/mux
// // net/http.NewServeMux()
// // https://github.com/gin-gonic/gin | still thinking...
// // https://github.com/labstack/echo | still thinking...

// // DRAFT: create router that can be adapt to above lib

// type IRouter interface {
// 	GET(path string, h func(context.Context, *http.Request) (any, error))
// 	POST(path string, h func(context.Context, *http.Request) (any, error))
// 	PUT(path string, h func(context.Context, *http.Request) (any, error))
// 	PATCH(path string, h func(context.Context, *http.Request) (any, error))
// 	DELETE(path string, h func(context.Context, *http.Request) (any, error))
// 	ServeHTTP(w http.ResponseWriter, r *http.Request)
// }

// type Adapter int

// const (
// 	TypeDefault Adapter = iota
// 	TypeHTTPRouter
// 	TypeChi
// 	TypeGorillaMux
// )

// type s struct {
// 	method, path string
// 	handler      func(context.Context, *http.Request) (any, error)
// }

// type Router struct {
// 	httpRouter       *httprouter.Router
// 	chiRouter        *chi.Mux
// 	gorillamuxRouter *mux.Router
// 	muxRouter        *http.ServeMux

// 	notFound                func(w http.ResponseWriter, r *http.Request)
// 	methodNotAllowed        func(w http.ResponseWriter, r *http.Request)
// 	adapter                 Adapter
// 	data                    []s
// 	isBeforeServeHTTPCalled bool
// }

// func New(a Adapter) *Router {
// 	r := &Router{data: make([]s, 0), adapter: a, isBeforeServeHTTPCalled: false}

// 	switch a {
// 	case TypeChi:
// 		r.chiRouter = chi.NewRouter()
// 	case TypeGorillaMux:
// 		r.gorillamuxRouter = mux.NewRouter()
// 	case TypeHTTPRouter:
// 		r.httpRouter = httprouter.New()
// 	default:
// 		r.muxRouter = http.NewServeMux()
// 	}

// 	return r
// }

// func (r *Router) GET(path string, h func(context.Context, *http.Request) (any, error)) {
// 	r.data = append(r.data, s{method: http.MethodGet, path: path, handler: h})
// }

// func (r *Router) POST(path string, h func(context.Context, *http.Request) (any, error)) {
// 	r.data = append(r.data, s{method: http.MethodPost, path: path, handler: h})
// }

// func (r *Router) PUT(path string, h func(context.Context, *http.Request) (any, error)) {
// 	r.data = append(r.data, s{method: http.MethodPut, path: path, handler: h})
// }

// func (r *Router) PATCH(path string, h func(context.Context, *http.Request) (any, error)) {
// 	r.data = append(r.data, s{method: http.MethodPatch, path: path, handler: h})
// }

// func (r *Router) DELETE(path string, h func(context.Context, *http.Request) (any, error)) {
// 	r.data = append(r.data, s{method: http.MethodDelete, path: path, handler: h})
// }

// func (r *Router) SetNotFound(f func(http.ResponseWriter, *http.Request)) {
// 	r.notFound = f
// }

// func (r *Router) SetMethodNotAllowed(f func(http.ResponseWriter, *http.Request)) {
// 	r.methodNotAllowed = f
// }

// func (r *Router) BeforeServeHTTP() {
// 	// register data path method and handler to active router base on adapter type
// 	// set true isBeforeServeHTTPCalled

// 	if r.notFound != nil {
// 		switch r.adapter {
// 		case TypeChi:
// 			r.chiRouter.NotFound(r.notFound)
// 		case TypeGorillaMux:
// 			r.gorillamuxRouter.NotFoundHandler = http.HandlerFunc(r.notFound)
// 		case TypeHTTPRouter:
// 			r.httpRouter.NotFound = http.HandlerFunc(r.notFound)
// 		default:
// 			r.muxRouter.Handle("/", http.HandlerFunc(r.notFound))
// 		}
// 	}

// 	if r.methodNotAllowed != nil {
// 		switch r.adapter {
// 		case TypeChi:
// 			r.chiRouter.MethodNotAllowed(r.methodNotAllowed)
// 		case TypeGorillaMux:
// 			r.gorillamuxRouter.MethodNotAllowedHandler = http.HandlerFunc(r.methodNotAllowed)
// 		case TypeHTTPRouter:
// 			r.httpRouter.MethodNotAllowed = http.HandlerFunc(r.methodNotAllowed)
// 		}
// 	}

// 	r.isBeforeServeHTTPCalled = true
// }

// func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	if !r.isBeforeServeHTTPCalled {
// 		panic("router: ServeHTTP call before BeforeServeHTTP")
// 	}

// 	switch {
// 	case r.chiRouter != nil:
// 		r.chiRouter.ServeHTTP(w, req)
// 	case r.gorillamuxRouter != nil:
// 		r.gorillamuxRouter.ServeHTTP(w, req)
// 	case r.httpRouter != nil:
// 		r.httpRouter.ServeHTTP(w, req)
// 	default:
// 		r.muxRouter.ServeHTTP(w, req)
// 	}
// }
