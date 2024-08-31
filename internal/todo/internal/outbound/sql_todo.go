package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/entity"
	"github.com/shandysiswandi/gostarter/pkg/config"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
)

type SQLTodo struct {
	db     *sql.DB
	qu     goqu.DialectWrapper
	config config.Config
}

func NewSQLTodo(db *sql.DB, config config.Config) *SQLTodo {
	qu := goqu.Dialect("mysql")
	if config.GetString("database.driver") == "postgres" {
		qu = goqu.Dialect("postgres")
	}

	return &SQLTodo{
		db:     db,
		qu:     qu,
		config: config,
	}
}

func (st *SQLTodo) Create(ctx context.Context, todo entity.Todo) error {
	query := func() (string, []any, error) {
		return st.qu.Insert("todos").
			Cols("id", "title", "description", "status").
			Vals([]any{todo.ID, todo.Title, todo.Description, todo.Status.String()}).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return entity.ErrTodoNotCreated
	}

	return err
}

func (st *SQLTodo) Delete(ctx context.Context, id uint64) error {
	query := func() (string, []any, error) {
		return st.qu.Delete("todos").Where(goqu.Ex{"id": id}).ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) GetByID(ctx context.Context, id uint64) (*entity.Todo, error) {
	query := func() (string, []any, error) {
		return st.qu.
			Select("id", "title", "description", "status").
			From("todos").
			Where(goqu.Ex{"id": id}).
			ToSQL()
	}

	todo, err := dbops.SQLGet[entity.Todo](ctx, st.db, query)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (st *SQLTodo) GetWithFilter(ctx context.Context, _ map[string]string) ([]entity.Todo, error) {
	query := func() (string, []any, error) {
		return st.qu.
			Select("id", "title", "description", "status").
			From("todos").
			ToSQL()
	}

	todos, err := dbops.SQLGets[entity.Todo](ctx, st.db, query)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (st *SQLTodo) UpdateStatus(ctx context.Context, id uint64, sts entity.TodoStatus) error {
	query := func() (string, []any, error) {
		return st.qu.
			Update("todos").
			Set(map[string]any{"status": sts}).
			Where(goqu.Ex{"id": id}).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) Update(ctx context.Context, todo entity.Todo) error {
	query := func() (string, []any, error) {
		return st.qu.
			Update("todos").
			Set(map[string]any{
				"title":       todo.Title,
				"description": todo.Description,
				"status":      todo.Status.String(),
			}).
			Where(goqu.Ex{"id": todo.ID}).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}
