package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type SQLRBAC struct {
	db        *sql.DB
	qu        goqu.DialectWrapper
	telemetry *telemetry.Telemetry
}

func NewSQLRBAC(db *sql.DB, qu goqu.DialectWrapper, tel *telemetry.Telemetry) *SQLRBAC {
	return &SQLRBAC{
		db:        db,
		qu:        qu,
		telemetry: tel,
	}
}

func (st *SQLRBAC) SaveRole(ctx context.Context, r domain.Role) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.SaveRole")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Insert("roles").
			Cols("id", "name", "description").
			Vals([]any{r.ID, r.Name, r.Description}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrRoleNotCreated
	}

	return err
}

func (st *SQLRBAC) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FindRoleByName")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "name", "description").
			From("roles").
			Where(goqu.Ex{"name": name}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Role](ctx, st.db, query)
}
