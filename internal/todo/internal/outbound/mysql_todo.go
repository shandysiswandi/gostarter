package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
)

type MysqlTodo struct {
	db *sql.DB
}

func NewMysqlTodo(db *sql.DB) *MysqlTodo {
	return &MysqlTodo{
		db: db,
	}
}

func (mt *MysqlTodo) Create(ctx context.Context, todo entity.Todo) error {
	query := func() (string, []any, error) {
		query := `INSERT INTO todos(id,title,description,status) VALUES(?,?,?,?)`

		return query, []any{todo.ID, todo.Title, todo.Description, todo.Status.String()}, nil
	}

	err := dbops.Exec(ctx, mt.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return entity.ErrTodoNotCreated
	}

	return err
}

func (mt *MysqlTodo) Delete(ctx context.Context, id uint64) error {
	query := func() (string, []any, error) {
		query := `DELETE FROM todos WHERE id=?`

		return query, []any{id}, nil
	}

	return dbops.Exec(ctx, mt.db, query)
}

func (mt *MysqlTodo) GetByID(ctx context.Context, id uint64) (*entity.Todo, error) {
	query := func() (string, []any, error) {
		query := "SELECT id, title, description,status FROM todos WHERE id = ?"

		return query, []any{id}, nil
	}

	todo, err := dbops.SQLGet[entity.Todo](ctx, mt.db, query)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (mt *MysqlTodo) GetWithFilter(ctx context.Context, _ map[string]string) ([]entity.Todo, error) {
	query := func() (string, []any, error) {
		return "SELECT id, title, description, status FROM todos", nil, nil
	}

	todos, err := dbops.SQLGets[entity.Todo](ctx, mt.db, query)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (mt *MysqlTodo) UpdateStatus(ctx context.Context, id uint64, sts entity.TodoStatus) error {
	query := func() (string, []any, error) {
		query := `UPDATE todos SET status=? WHERE id=?`

		return query, []any{sts.String(), id}, nil
	}

	return dbops.Exec(ctx, mt.db, query)
}

func (mt *MysqlTodo) Update(ctx context.Context, todo entity.Todo) error {
	query := func() (string, []any, error) {
		query := `UPDATE todos SET title=?,description=?,status=? WHERE id=?`

		return query, []any{todo.Title, todo.Description, todo.Status.String(), todo.ID}, nil
	}

	return dbops.Exec(ctx, mt.db, query)
}
