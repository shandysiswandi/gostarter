package http

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/region/internal/usecase"
	pkghttp "github.com/shandysiswandi/gostarter/pkg/http"
)

func RegisterRESTEndpoint(router *httprouter.Router, h *Endpoint) {
	serve := pkghttp.NewServe(
		pkghttp.WithMiddlewares(pkghttp.Recovery),
	)

	router.Handler(http.MethodGet, "/regions/search", serve.Endpoint(h.Search))
}

type Endpoint struct {
	SearchUC usecase.Search
}

func (e *Endpoint) Search(ctx context.Context, r *http.Request) (any, error) {
	by := r.URL.Query().Get("by")
	pid := r.URL.Query().Get("pid")
	ids := r.URL.Query().Get("ids")

	resp, err := e.SearchUC.Execute(ctx, usecase.SearchInput{By: by, ParentID: pid, IDs: ids})
	if err != nil {
		return nil, err
	}

	return SearchResponse{Type: by, Regions: FromListEntityRegion(resp.Regions)}, nil
}
