package inbound

import (
	"encoding/json"
	"strconv"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type httpEndpoint struct {
	tel *telemetry.Telemetry

	createUC       domain.Create
	deleteUC       domain.Delete
	findUC         domain.Find
	fetchUC        domain.Fetch
	updateStatusUC domain.UpdateStatus
	updateUC       domain.Update
}

func (h *httpEndpoint) Create(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "todo.inbound.httpEndpoint.Create")
	defer span.End()

	var req CreateRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.createUC.Call(ctx, domain.CreateInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return CreateResponse{ID: resp.ID}, nil
}

func (h *httpEndpoint) Delete(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "todo.inbound.httpEndpoint.Delete")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := h.deleteUC.Call(ctx, domain.DeleteInput{ID: id})
	if err != nil {
		return nil, err
	}

	return DeleteResponse{ID: resp.ID}, nil
}

func (h *httpEndpoint) Find(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "todo.inbound.httpEndpoint.Find")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := h.findUC.Call(ctx, domain.FindInput{ID: id})
	if err != nil {
		return nil, err
	}

	return FindResponse{
		ID:          resp.ID,
		UserID:      resp.UserID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      resp.Status.String(),
	}, nil
}

func (h *httpEndpoint) Fetch(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "todo.inbound.httpEndpoint.Fetch")
	defer span.End()

	resp, err := h.fetchUC.Call(ctx, domain.FetchInput{
		Cursor: c.Query("cursor"),
		Limit:  c.Query("limit"),
		Status: c.Query("status"),
	})
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, 0)
	for _, todo := range resp.Todos {
		todos = append(todos, Todo{
			ID:          todo.ID,
			UserID:      todo.UserID,
			Title:       todo.Title,
			Description: todo.Description,
			Status:      todo.Status.String(),
		})
	}

	return FetchResponse{
		Todos: todos,
		Pagination: Pagination{
			NextCursor: resp.NextCursor,
			HasMore:    resp.HasMore,
		},
	}, nil
}

func (h *httpEndpoint) UpdateStatus(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "todo.inbound.httpEndpoint.UpdateStatus")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	var req UpdateStatusRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.updateStatusUC.Call(ctx, domain.UpdateStatusInput{
		ID:     id,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return UpdateStatusResponse{
		ID:     id,
		Status: resp.Status.String(),
	}, nil
}

func (h *httpEndpoint) Update(c framework.Context) (any, error) {
	ctx, span := h.tel.Tracer().Start(c.Context(), "todo.inbound.httpEndpoint.Update")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	var req UpdateRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.updateUC.Call(ctx, domain.UpdateInput{
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
		UserID:      resp.UserID,
		Title:       resp.Title,
		Description: resp.Description,
		Status:      resp.Status.String(),
	}, nil
}
