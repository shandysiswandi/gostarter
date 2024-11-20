package inbound

import (
	"encoding/json"
	"strconv"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework/httpserver"
)

type httpEndpoint struct {
	createUC       domain.Create
	deleteUC       domain.Delete
	findUC         domain.Find
	fetchUC        domain.Fetch
	updateStatusUC domain.UpdateStatus
	updateUC       domain.Update
}

func (e *httpEndpoint) Create(c httpserver.Context) (any, error) {
	var req CreateRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.createUC.Call(c.Context(), domain.CreateInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return CreateResponse{ID: resp.ID}, nil
}

func (e *httpEndpoint) Delete(c httpserver.Context) (any, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.deleteUC.Call(c.Context(), domain.DeleteInput{ID: id})
	if err != nil {
		return nil, err
	}

	return DeleteResponse{ID: resp.ID}, nil
}

func (e *httpEndpoint) Find(c httpserver.Context) (any, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.findUC.Call(c.Context(), domain.FindInput{ID: id})
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

func (e *httpEndpoint) Fetch(c httpserver.Context) (any, error) {
	resp, err := e.fetchUC.Call(c.Context(), domain.FetchInput{
		ID:          c.Query("id"),
		Title:       c.Query("title"),
		Description: c.Query("description"),
		Status:      c.Query("status"),
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

func (e *httpEndpoint) UpdateStatus(c httpserver.Context) (any, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.updateStatusUC.Call(c.Context(), domain.UpdateStatusInput{ID: id, Status: req.Status})
	if err != nil {
		return nil, err
	}

	return UpdateStatusResponse{ID: id, Status: resp.Status.String()}, nil
}

func (e *httpEndpoint) Update(c httpserver.Context) (any, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	var req UpdateRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := e.updateUC.Call(c.Context(), domain.UpdateInput{
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
