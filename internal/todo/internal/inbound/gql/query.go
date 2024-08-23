package gql

import (
	"context"
	"strconv"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/usecase"
)

type query struct{ *Endpoint }

func (q *query) GetWithFilter(ctx context.Context, in *ql.FilterInput) ([]ql.Todo, error) {
	input := usecase.GetWithFilterInput{}
	if in != nil {
		input = usecase.GetWithFilterInput{
			ID:          getString(in.ID),
			Title:       getString(in.Title),
			Description: getString(in.Description),
			Status:      getStatusString(in.Status),
		}
	}

	resp, err := q.GetWithFilterUC.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	todos := make([]ql.Todo, 0)
	for _, todo := range resp.Todos {
		todos = append(todos, ql.Todo{
			ID:          strconv.FormatUint(todo.ID, 10),
			Title:       todo.Title,
			Description: todo.Description,
			Status:      ql.Status(todo.Status.String()),
		})
	}

	return todos, nil
}

func (q *query) GetByID(ctx context.Context, id string) (*ql.Todo, error) {
	idu64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errfailedParseToUint
	}

	resp, err := q.GetByIDUC.Execute(ctx, usecase.GetByIDInput{ID: idu64})
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
