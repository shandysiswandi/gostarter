package gql

import (
	"context"
	"strconv"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
)

var errFailedParseToUint = goerror.NewInvalidInput("failed parse id to uint", nil)

func getString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}

	return ""
}

func getStatusString(status *ql.Status) string {
	if status != nil && status.IsValid() {
		return status.String()
	}

	return ""
}

type Endpoint struct {
	ql.Resolver

	FindUC         domain.Find
	FetchUC        domain.Fetch
	CreateUC       domain.Create
	DeleteUC       domain.Delete
	UpdateUC       domain.Update
	UpdateStatusUC domain.UpdateStatus
}

func (e *Endpoint) Mutation() ql.MutationResolver { return e }

func (e *Endpoint) Query() ql.QueryResolver { return e }

func (e *Endpoint) Fetch(ctx context.Context, in *ql.FetchInput) ([]ql.Todo, error) {
	input := domain.FetchInput{}
	if in != nil {
		input = domain.FetchInput{
			ID:          getString(in.ID),
			Title:       getString(in.Title),
			Description: getString(in.Description),
			Status:      getStatusString(in.Status),
		}
	}

	resp, err := e.FetchUC.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	todos := make([]ql.Todo, 0)
	for _, todo := range resp {
		todos = append(todos, ql.Todo{
			ID:          strconv.FormatUint(todo.ID, 10),
			Title:       todo.Title,
			Description: todo.Description,
			Status:      ql.Status(todo.Status.String()),
		})
	}

	return todos, nil
}

func (e *Endpoint) Find(ctx context.Context, id string) (*ql.Todo, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.FindUC.Execute(ctx, domain.FindInput{ID: idu64})
	if err != nil {
		return nil, err
	}

	return &ql.Todo{
		ID:          strconv.FormatUint(resp.ID, 10),
		Title:       resp.Title,
		Description: resp.Description,
		Status:      ql.Status(resp.Status.String()),
	}, nil
}

func (e *Endpoint) Create(ctx context.Context, in ql.CreateInput) (*ql.Todo, error) {
	resp, err := e.CreateUC.Execute(ctx, domain.CreateInput{Title: in.Title, Description: in.Description})
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

func (e *Endpoint) Delete(ctx context.Context, id string) (string, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", errFailedParseToUint
	}

	resp, err := e.DeleteUC.Execute(ctx, domain.DeleteInput{ID: idu64})
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(resp.ID, 10), nil
}

func (e *Endpoint) UpdateStatus(ctx context.Context, id string, status ql.Status) (
	*ql.UpdateStatusResponse, error,
) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.UpdateStatusUC.Execute(ctx, domain.UpdateStatusInput{ID: idu64, Status: status.String()})
	if err != nil {
		return nil, err
	}

	return &ql.UpdateStatusResponse{
		ID:     strconv.FormatUint(resp.ID, 10),
		Status: ql.Status(resp.Status.String()),
	}, nil
}

func (e *Endpoint) Update(ctx context.Context, id string, in ql.UpdateInput) (*ql.UpdateResponse, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.UpdateUC.Execute(ctx, domain.UpdateInput{
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
