package inbound

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	pkghttp "github.com/shandysiswandi/gostarter/pkg/http"
)

var ErrDecodeBody = errors.New("failed to decode body")

func RegisterHTTP(router *httprouter.Router, h *Endpoint) {
	serve := pkghttp.NewServe(
		pkghttp.WithMiddlewares(pkghttp.Recovery),
	)

	router.Handler(http.MethodGet, "/shortly/:key", serve.Endpoint(h.Get))
	router.Handler(http.MethodPost, "/shortly", serve.Endpoint(h.Set))
}

type Endpoint struct {
	GetUC domain.Get
	SetUC domain.Set
}

func (e *Endpoint) Get(ctx context.Context, _ *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(ctx)
	key := params.ByName("key")

	resp, err := e.GetUC.Call(ctx, domain.GetInput{Key: key})
	if err != nil {
		return nil, err
	}

	if resp.URL == "" {
		return map[string]any{"url": "NOT_FOUND"}, nil
	}

	return map[string]any{"url": resp.URL}, nil
}

func (e *Endpoint) Set(ctx context.Context, r *http.Request) (any, error) {
	var req struct {
		URL  string `json:"url"`
		Slug string `json:"slug"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, ErrDecodeBody
	}

	resp, err := e.SetUC.Call(ctx, domain.SetInput{URL: req.URL, Slug: req.Slug})
	if err != nil {
		return nil, err
	}

	return map[string]any{"key": resp.Key}, nil
}
