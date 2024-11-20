package httpserver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Context interface {
	Context() context.Context

	Request() *http.Request

	Header() http.Header

	Body() io.ReadCloser

	Query(key string) string

	Queries(key string) []string

	Param(key string) string
}

type routerCtx struct {
	r *http.Request
}

func (rc *routerCtx) Context() context.Context {
	return rc.r.Context()
}

func (rc *routerCtx) Request() *http.Request {
	return rc.r
}

func (rc *routerCtx) Header() http.Header {
	return rc.r.Header
}

func (rc *routerCtx) Body() io.ReadCloser {
	return rc.r.Body
}

func (rc *routerCtx) Query(key string) string {
	return rc.r.URL.Query().Get(key)
}

func (rc *routerCtx) Queries(key string) []string {
	return rc.r.URL.Query()[key]
}

func (rc *routerCtx) Param(key string) string {
	return httprouter.ParamsFromContext(rc.Context()).ByName(key)
}

type TestCtx struct {
	r       *http.Request
	param   map[string]string
	queries map[string][]string
	headers map[string][]string
	mu      sync.RWMutex
}

func NewTestContext(method string, target string, body io.Reader) *TestCtx {
	r := httptest.NewRequest(method, target, body)

	return &TestCtx{
		r:       r,
		param:   make(map[string]string),
		queries: make(map[string][]string),
		headers: make(map[string][]string),
	}
}

func (tc *TestCtx) SetParam(key, value string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.param[key] = value
}

func (tc *TestCtx) SetQuery(key string, values ...string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.queries[key] = append(tc.queries[key], values...)
}

func (tc *TestCtx) SetHeader(key string, values ...string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.headers[key] = append(tc.headers[key], values...)
}

func (tc *TestCtx) Build() *routerCtx {
	q := tc.r.URL.Query()
	for key, query := range tc.queries {
		for _, value := range query {
			q.Add(key, value)
		}
	}
	tc.r.URL.RawQuery = q.Encode()

	for key, header := range tc.headers {
		for _, value := range header {
			tc.r.Header.Add(key, value)
		}
	}

	var params httprouter.Params
	for key, value := range tc.param {
		params = append(params, httprouter.Param{
			Key:   key,
			Value: value,
		})
	}
	ctx := context.WithValue(tc.r.Context(), httprouter.ParamsKey, params)

	return &routerCtx{r: tc.r.WithContext(ctx)}
}
