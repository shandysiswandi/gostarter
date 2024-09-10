package gql

import (
	"context"
	"strconv"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
)

type query struct{ *Endpoint }

func (q *query) Fetch(ctx context.Context, in *ql.FetchInput) ([]ql.Todo, error) {
	input := domain.FetchInput{}
	if in != nil {
		input = domain.FetchInput{
			ID:          getString(in.ID),
			Title:       getString(in.Title),
			Description: getString(in.Description),
			Status:      getStatusString(in.Status),
		}
	}

	resp, err := q.FetchUC.Execute(ctx, input)
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

func (q *query) Find(ctx context.Context, id string) (*ql.Todo, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := q.FindUC.Execute(ctx, domain.FindInput{ID: idu64})
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
