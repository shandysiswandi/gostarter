package http

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/region/internal/usecase"
	"github.com/shandysiswandi/gostarter/pkg/http/middleware"
	"github.com/shandysiswandi/gostarter/pkg/http/serve"
)

func RegisterRESTEndpoint(router *httprouter.Router, h *Endpoint) {
	serve := serve.New(serve.WithMiddlewares(middleware.Recovery))

	router.Handler(http.MethodGet, "/regions/search", serve.Endpoint(h.Search))
}

type Endpoint struct {
	search usecase.Search
}

func NewEndpoint(search usecase.Search) *Endpoint {
	return &Endpoint{search: search}
}

func (e *Endpoint) Search(ctx context.Context, r *http.Request) (any, error) {
	by := r.URL.Query().Get("by")
	pid := r.URL.Query().Get("pid")
	ids := r.URL.Query().Get("ids")

	resp, err := e.search.Execute(ctx, usecase.SearchInput{By: by, ParentID: pid, IDs: ids})
	if err != nil {
		return nil, err
	}

	return SearchResponse{Type: by, Regions: FromListEntityRegion(resp.Regions)}, nil
}
