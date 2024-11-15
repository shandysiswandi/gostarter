package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
)

var errFailedParseToUint = goerror.NewInvalidInput("failed parse id to uint", nil)

type Endpoint struct {
	createUC       domain.Create
	deleteUC       domain.Delete
	findUC         domain.Find
	fetchUC        domain.Fetch
	updateStatusUC domain.UpdateStatus
	updateUC       domain.Update
}

func NewEndpoint(
	createUC domain.Create,
	deleteUC domain.Delete,
	findUC domain.Find,
	fetchUC domain.Fetch,
	updateStatusUC domain.UpdateStatus,
	updateUC domain.Update,
) *Endpoint {
	return &Endpoint{
		createUC:       createUC,
		deleteUC:       deleteUC,
		findUC:         findUC,
		fetchUC:        fetchUC,
		updateStatusUC: updateStatusUC,
		updateUC:       updateUC,
	}
}

func (e *Endpoint) Create(ctx context.Context, r *http.Request) (any, error) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	resp, err := e.createUC.Execute(ctx, domain.CreateInput{Title: req.Title, Description: req.Description})
	if err != nil {
		return nil, err
	}

	return CreateResponse{ID: resp.ID}, nil
}

func (e *Endpoint) Delete(ctx context.Context, _ *http.Request) (any, error) {
	idStr := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.deleteUC.Execute(ctx, domain.DeleteInput{ID: id})
	if err != nil {
		return nil, err
	}

	return DeleteResponse{ID: resp.ID}, nil
}

func (e *Endpoint) Find(ctx context.Context, _ *http.Request) (any, error) {
	idStr := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.findUC.Execute(ctx, domain.FindInput{ID: id})
	if err != nil {
		return nil, err
	}

	return FindResponse{
		ID:          resp.ID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      resp.Status.String(),
	}, nil
}

func (e *Endpoint) Fetch(ctx context.Context, r *http.Request) (any, error) {
	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	description := r.URL.Query().Get("description")
	status := r.URL.Query().Get("status")

	resp, err := e.fetchUC.Execute(ctx, domain.FetchInput{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, 0)
	for _, todo := range resp {
		todos = append(todos, Todo{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status.String(),
		})
	}

	return FetchResponse{Todos: todos}, nil
}

func (e *Endpoint) UpdateStatus(ctx context.Context, r *http.Request) (any, error) {
	var req UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	idStr := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.updateStatusUC.Execute(ctx, domain.UpdateStatusInput{ID: id, Status: req.Status})
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

	idStr := httprouter.ParamsFromContext(ctx).ByName("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.updateUC.Execute(ctx, domain.UpdateInput{
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
