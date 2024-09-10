package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
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

func (st *SQLTodo) Create(ctx context.Context, todo domain.Todo) error {
	query := func() (string, []any, error) {
		return st.qu.Insert("todos").Cols("id", "title", "description", "status").
			Vals([]any{todo.ID, todo.Title, todo.Description, todo.Status}).
			Prepared(true).ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrTodoNotCreated
	}

	return err
}

func (st *SQLTodo) Delete(ctx context.Context, id uint64) error {
	query := func() (string, []any, error) {
		return st.qu.Delete("todos").Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) Find(ctx context.Context, id uint64) (*domain.Todo, error) {
	query := func() (string, []any, error) {
		return st.qu.Select("id", "title", "description", "status").
			From("todos").Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	}

	return dbops.SQLGet[domain.Todo](ctx, st.db, query)
}

func (st *SQLTodo) Fetch(ctx context.Context, _ map[string]string) ([]domain.Todo, error) {
	query := func() (string, []any, error) {
		return st.qu.Select("id", "title", "description", "status").
			From("todos").Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.Todo](ctx, st.db, query)
}

func (st *SQLTodo) UpdateStatus(ctx context.Context, id uint64, sts domain.TodoStatus) error {
	query := func() (string, []any, error) {
		return st.qu.Update("todos").Set(map[string]any{"status": sts}).
			Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) Update(ctx context.Context, todo domain.Todo) error {
	query := func() (string, []any, error) {
		return st.qu.Update("todos").Set(map[string]any{
			"title":       todo.Title,
			"description": todo.Description,
			"status":      todo.Status,
		}).Where(goqu.Ex{"id": todo.ID}).Prepared(true).ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}
