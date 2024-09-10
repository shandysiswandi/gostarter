package gql

import (
	"context"
	"strconv"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type mutation struct{ *Endpoint }

func (m *mutation) Create(ctx context.Context, in ql.CreateInput) (*ql.Todo, error) {
	resp, err := m.CreateUC.Execute(ctx, domain.CreateInput{Title: in.Title, Description: in.Description})
	if err != nil {
		return nil, err
	}

	return &ql.Todo{
		ID:          strconv.FormatUint(resp.ID, 10),
		Title:       in.Title,
		Description: in.Description,
		Status:      ql.Status(domain.TodoStatusInitiate.String()),
	}, nil
}

func (m *mutation) Delete(ctx context.Context, id string) (string, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", errFailedParseToUint
	}

	resp, err := m.DeleteUC.Execute(ctx, domain.DeleteInput{ID: idu64})
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(resp.ID, 10), nil
}

func (m *mutation) UpdateStatus(ctx context.Context, id string, status ql.Status) (*ql.UpdateStatusResponse, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := m.UpdateStatusUC.Execute(ctx, domain.UpdateStatusInput{ID: idu64, Status: status.String()})
	if err != nil {
		return nil, err
	}

	return &ql.UpdateStatusResponse{ID: strconv.FormatUint(resp.ID, 10), Status: ql.Status(resp.Status.String())}, nil
}

func (m *mutation) Update(ctx context.Context, id string, in ql.UpdateInput) (*ql.UpdateResponse, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := m.UpdateUC.Execute(ctx, domain.UpdateInput{
		ID:          idu64,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status.String(),
	})
	if err != nil {
		return nil, err
	}

	return &ql.UpdateResponse{
		ID:          strconv.FormatUint(resp.ID, 10),
		Title:       resp.Title,
		Description: resp.Description,
		Status:      ql.Status(resp.Status.String()),
	}, nil
}
