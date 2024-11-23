package inbound

import (
	"context"
	"strconv"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

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

type gqlEndpoint struct {
	ql.Resolver

	findUC         domain.Find
	fetchUC        domain.Fetch
	createUC       domain.Create
	deleteUC       domain.Delete
	updateUC       domain.Update
	updateStatusUC domain.UpdateStatus
}

func (e *gqlEndpoint) Mutation() ql.MutationResolver { return e }

func (e *gqlEndpoint) Query() ql.QueryResolver { return e }

func (e *gqlEndpoint) Fetch(ctx context.Context, in *ql.FetchInput) ([]ql.Todo, error) {
	input := domain.FetchInput{}
	if in != nil {
		input = domain.FetchInput{
			ID:          getString(in.ID),
			Title:       getString(in.Title),
			Description: getString(in.Description),
			Status:      getStatusString(in.Status),
		}
	}

	resp, err := e.fetchUC.Call(ctx, input)
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

func (e *gqlEndpoint) Find(ctx context.Context, id string) (*ql.Todo, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.findUC.Call(ctx, domain.FindInput{ID: idu64})
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

func (e *gqlEndpoint) Create(ctx context.Context, in ql.CreateInput) (string, error) {
	resp, err := e.createUC.Call(ctx, domain.CreateInput{Title: in.Title, Description: in.Description})
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(resp.ID, 10), nil
}

func (e *gqlEndpoint) Delete(ctx context.Context, id string) (string, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return "", errFailedParseToUint
	}

	resp, err := e.deleteUC.Call(ctx, domain.DeleteInput{ID: idu64})
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(resp.ID, 10), nil
}

func (e *gqlEndpoint) UpdateStatus(ctx context.Context, in ql.UpdateStatusInput) (
	*ql.UpdateStatusOutput, error,
) {
	idu64, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.updateStatusUC.Call(ctx, domain.UpdateStatusInput{ID: idu64, Status: in.Status.String()})
	if err != nil {
		return nil, err
	}

	return &ql.UpdateStatusOutput{
		ID:     strconv.FormatUint(resp.ID, 10),
		Status: ql.Status(resp.Status.String()),
	}, nil
}

func (e *gqlEndpoint) Update(ctx context.Context, in ql.UpdateInput) (*ql.Todo, error) {
	idu64, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := e.updateUC.Call(ctx, domain.UpdateInput{
		ID:          idu64,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status.String(),
	})
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
