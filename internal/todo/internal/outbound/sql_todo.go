package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/goreng/enum"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type SQLTodo struct {
	db  *sql.DB
	qu  goqu.DialectWrapper
	tel *telemetry.Telemetry
}

func NewSQLTodo(db *sql.DB, qu goqu.DialectWrapper, tel *telemetry.Telemetry) *SQLTodo {
	return &SQLTodo{
		db:  db,
		qu:  qu,
		tel: tel,
	}
}

func (st *SQLTodo) Create(ctx context.Context, todo domain.Todo) error {
	ctx, span := st.tel.Tracer().Start(ctx, "auth.outbound.SQLTodo.Create")
	defer span.End()

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
	ctx, span := st.tel.Tracer().Start(ctx, "auth.outbound.SQLTodo.Delete")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Delete("todos").
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLTodo) Find(ctx context.Context, id uint64) (*domain.Todo, error) {
	ctx, span := st.tel.Tracer().Start(ctx, "auth.outbound.SQLTodo.Find")
	defer span.End()

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
	ctx, span := st.tel.Tracer().Start(ctx, "auth.outbound.SQLTodo.Fetch")
	defer span.End()

	cursor, hasCursor := filter["cursor"].(uint64)
	limit, hasLimit := filter["limit"].(int)
	status, hasStatus := filter["status"].(enum.Enum[domain.TodoStatus])

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

func (st *SQLTodo) UpdateStatus(ctx context.Context, id uint64, sts enum.Enum[domain.TodoStatus]) error {
	ctx, span := st.tel.Tracer().Start(ctx, "auth.outbound.SQLTodo.UpdateStatus")
	defer span.End()

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
	ctx, span := st.tel.Tracer().Start(ctx, "auth.outbound.SQLTodo.Update")
	defer span.End()

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
