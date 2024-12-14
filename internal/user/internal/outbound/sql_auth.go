package outbound

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type SQLUser struct {
	db        *sql.DB
	qu        goqu.DialectWrapper
	telemetry *telemetry.Telemetry
}

func NewSQLUser(db *sql.DB, qu goqu.DialectWrapper, tel *telemetry.Telemetry) *SQLUser {
	return &SQLUser{
		db:        db,
		qu:        qu,
		telemetry: tel,
	}
}

func (st *SQLUser) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.SQLUser.FindUserByEmail")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "name", "email", "password").
			From("users").
			Where(goqu.Ex{"email": email}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.User](ctx, st.db, query)
}

func (st *SQLUser) Update(ctx context.Context, user map[string]any) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.Update")
	defer span.End()

	id := user["id"]
	delete(user, "id")

	query := func() (string, []any, error) {
		return st.qu.Update("users").
			Set(user).
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLUser) DeleteTokenByAccess(ctx context.Context, token string) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.DeleteTokenByAccess")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Delete("tokens").
			Where(goqu.Ex{"access_token": token}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}
