package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
)

type SQLTodo struct {
	db *sql.DB
	qu goqu.DialectWrapper
}

func NewSQLTodo(db *sql.DB, qu goqu.DialectWrapper) *SQLTodo {
	return &SQLTodo{
		db: db,
		qu: qu,
	}
}

func (st *SQLTodo) Create(ctx context.Context, todo domain.Todo) error {
	query := func() (string, []any, error) {
		return st.qu.Insert("todos").
			Cols("id", "user_id", "title", "description", "status").
			Vals([]any{todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrTodoNotCreated
	}

	return err
}

func (st *SQLTodo) Delete(ctx context.Context, id uint64) error {
	query := func() (string, []any, error) {
		return st.qu.Delete("todos").
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) Find(ctx context.Context, id uint64) (*domain.Todo, error) {
	query := func() (string, []any, error) {
		return st.qu.Select("id", "user_id", "title", "description", "status").
			From("todos").
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Todo](ctx, st.db, query)
}

func (st *SQLTodo) Fetch(ctx context.Context, filter map[string]any) ([]domain.Todo, error) {
	cursor, hasCursor := filter["cursor"].(uint64)
	limit, hasLimit := filter["limit"].(int)
	status, hasStatus := filter["status"].(domain.TodoStatus)

	query := func() (string, []any, error) {
		q := st.qu.Select("id", "user_id", "title", "description", "status").
			From("todos")

		if hasCursor && cursor > 0 {
			q = q.Where(goqu.Ex{"id": goqu.Op{"gt": cursor}})
		}

		if hasStatus {
			q = q.Where(goqu.Ex{"status": status})
		}

		if hasLimit {
			q = q.Limit(uint(limit + 1))
		}

		return q.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.Todo](ctx, st.db, query)
}

func (st *SQLTodo) UpdateStatus(ctx context.Context, id uint64, sts domain.TodoStatus) error {
	query := func() (string, []any, error) {
		return st.qu.Update("todos").
			Set(map[string]any{"status": sts}).
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) Update(ctx context.Context, todo domain.Todo) error {
	query := func() (string, []any, error) {
		return st.qu.Update("todos").
			Set(map[string]any{
				"user_id":     todo.UserID,
				"title":       todo.Title,
				"description": todo.Description,
				"status":      todo.Status,
			}).
			Where(goqu.Ex{"id": todo.ID}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}
