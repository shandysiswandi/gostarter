package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
)

type PostgresTodo struct {
	db *sql.DB
}

func NewPostgresTodo(db *sql.DB) *PostgresTodo {
	return &PostgresTodo{
		db: db,
	}
}

func (pt *PostgresTodo) Create(ctx context.Context, todo entity.Todo) error {
	query := func() (string, []any, error) {
		query := `INSERT INTO todos(id, title, description, status) VALUES($1, $2, $3, $4)`

		return query, []any{todo.ID, todo.Title, todo.Description, todo.Status.String()}, nil
	}

	err := dbops.Exec(ctx, pt.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return entity.ErrTodoNotCreated
	}

	return err
}

func (pt *PostgresTodo) Delete(ctx context.Context, id uint64) error {
	query := func() (string, []any, error) {
		query := `DELETE FROM todos WHERE id=$1`

		return query, []any{id}, nil
	}

	return dbops.Exec(ctx, pt.db, query)
}

func (pt *PostgresTodo) GetByID(ctx context.Context, id uint64) (*entity.Todo, error) {
	query := func() (string, []any, error) {
		query := "SELECT id, title, description, status FROM todos WHERE id = $1"

		return query, []any{id}, nil
	}

	todo, err := dbops.SQLGet[entity.Todo](ctx, pt.db, query)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (pt *PostgresTodo) GetWithFilter(ctx context.Context, _ map[string]string) ([]entity.Todo, error) {
	query := func() (string, []any, error) {
		return "SELECT id, title, description, status FROM todos", nil, nil
	}

	todos, err := dbops.SQLGets[entity.Todo](ctx, pt.db, query)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (pt *PostgresTodo) UpdateStatus(ctx context.Context, id uint64, sts entity.TodoStatus) error {
	query := func() (string, []any, error) {
		query := `UPDATE todos SET status=$1 WHERE id=$2`

		return query, []any{sts.String(), id}, nil
	}

	return dbops.Exec(ctx, pt.db, query)
}

func (pt *PostgresTodo) Update(ctx context.Context, todo entity.Todo) error {
	query := func() (string, []any, error) {
		query := `UPDATE todos SET title=$1, description=$2, status=$3 WHERE id=$4`

		return query, []any{todo.Title, todo.Description, todo.Status.String(), todo.ID}, nil
	}

	return dbops.Exec(ctx, pt.db, query)
}
