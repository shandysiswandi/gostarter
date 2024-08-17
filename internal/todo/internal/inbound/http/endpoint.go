package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
	pkghttp "github.com/shandysiswandi/gostarter/pkg/http"
)

func RegisterRESTEndpoint(router *httprouter.Router, h *Endpoint) {
	serve := pkghttp.NewServe(
		pkghttp.WithMiddlewares(recovery),
	)

	router.Handler(http.MethodGet, "/todos/:id", serve.Endpoint(h.GetByID))
	router.Handler(http.MethodGet, "/todos", serve.Endpoint(h.GetWithFilter))
	router.Handler(http.MethodPost, "/todos", serve.Endpoint(h.Create))
	router.Handler(http.MethodPut, "/todos/:id", serve.Endpoint(h.Update))
	router.Handler(http.MethodPatch, "/todos/:id/status", serve.Endpoint(h.UpdateStatus))
	router.Handler(http.MethodDelete, "/todos/:id", serve.Endpoint(h.Delete))
}

type Endpoint struct {
	GetByIDUC       usecase.GetByID
	GetWithFilterUC usecase.GetWithFilter
	CreateUC        usecase.Create
	DeleteUC        usecase.Delete
	UpdateUC        usecase.Update
	UpdateStatusUC  usecase.UpdateStatus
}

func (e *Endpoint) Create(ctx context.Context, r *http.Request) (any, error) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	resp, err := e.CreateUC.Execute(ctx, usecase.CreateInput{Title: req.Title, Description: req.Description})
	if err != nil {
		return nil, err
	}

	return CreateResponse{ID: resp.ID}, nil
}

func (e *Endpoint) Delete(ctx context.Context, _ *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(ctx)
	idstr := params.ByName("id")

	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return nil, err
	}

	resp, err := e.DeleteUC.Execute(ctx, usecase.DeleteInput{ID: id})
	if err != nil {
		return nil, err
	}

	return DeleteResponse{ID: resp.ID}, nil
}

func (e *Endpoint) GetByID(ctx context.Context, _ *http.Request) (any, error) {
	params := httprouter.ParamsFromContext(ctx)
	idstr := params.ByName("id")

	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return nil, err
	}

	resp, err := e.GetByIDUC.Execute(ctx, usecase.GetByIDInput{ID: id})
	if err != nil {
		return nil, err
	}

	return GetByIDResponse{
		ID:          resp.ID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      resp.Status.String(),
	}, nil
}

func (e *Endpoint) GetWithFilter(ctx context.Context, r *http.Request) (any, error) {
	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")
	status := r.URL.Query().Get("status")

	resp, err := e.GetWithFilterUC.Execute(ctx, usecase.GetWithFilterInput{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, 0)
	for _, todo := range resp.Todos {
		todos = append(todos, Todo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status.String(),
		})
	}

	return GetWithFilterResponse{
		Todos: todos,
	}, nil
}

func (e *Endpoint) UpdateStatus(ctx context.Context, r *http.Request) (any, error) {
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	params := httprouter.ParamsFromContext(ctx)
	idstr := params.ByName("id")

	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return nil, err
	}

	resp, err := e.UpdateStatusUC.Execute(ctx, usecase.UpdateStatusInput{ID: id, Status: req.Status})
	if err != nil {
		return nil, err
	}

	return UpdateStatusResponse{ID: id, Status: resp.Status.String()}, nil
}

func (e *Endpoint) Update(ctx context.Context, r *http.Request) (any, error) {
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	params := httprouter.ParamsFromContext(ctx)
	idstr := params.ByName("id")

	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return nil, err
	}

	resp, err := e.UpdateUC.Execute(ctx, usecase.UpdateInput{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}

	return UpdateResponse{
		ID:          id,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      resp.Status.String(),
	}, nil
}
