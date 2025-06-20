package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
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

func (sr *SQLRBAC) SaveRole(ctx context.Context, r domain.Role) error {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.SaveRole")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Insert("roles").
			Cols("id", "name", "description").
			Vals([]any{r.ID, r.Name, r.Description}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, sr.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrRoleNotCreated
	}

	return err
}

func (sr *SQLRBAC) EditRole(ctx context.Context, r domain.Role) error {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.EditRole")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Update("roles").
			Set(map[string]any{
				"name":        r.Name,
				"description": r.Description,
			}).
			Where(goqu.Ex{"id": r.ID}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, sr.db, query)
}

func (sr *SQLRBAC) FindRole(ctx context.Context, id uint64) (*domain.Role, error) {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FindRole")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Select("id", "name", "description").
			From("roles").
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Role](ctx, sr.db, query)
}

func (sr *SQLRBAC) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FindRoleByName")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Select("id", "name", "description").
			From("roles").
			Where(goqu.Ex{"name": name}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Role](ctx, sr.db, query)
}

//nolint:dupl // #11 this is not duplicate with #12
func (sr *SQLRBAC) FetchRole(ctx context.Context, filter map[string]any) ([]domain.Role, error) {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FetchRole")
	defer span.End()

	cursor, hasCursor := filter["cursor"].(uint64)
	limit, hasLimit := filter["limit"].(int)
	name, hasName := filter["name"].(string)

	query := func() (string, []any, error) {
		q := sr.qu.Select("id", "name", "description").
			From("roles")

		if hasCursor && cursor > 0 {
			q = q.Where(goqu.Ex{"id": goqu.Op{"gt": cursor}})
		}

		if hasName {
			q = q.Where(goqu.Ex{"name": goqu.Op{"like": "%" + name + "%"}})
		}

		if hasLimit {
			q = q.Limit(uint(limit + 1))
		}

		return q.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.Role](ctx, sr.db, query)
}

func (sr *SQLRBAC) SavePermission(ctx context.Context, p domain.Permission) error {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.SavePermission")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Insert("permissions").
			Cols("id", "name", "description").
			Vals([]any{p.ID, p.Name, p.Description}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, sr.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrPermissionNotCreated
	}

	return err
}

func (sr *SQLRBAC) EditPermission(ctx context.Context, p domain.Permission) error {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.EditPermission")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Update("permissions").
			Set(map[string]any{
				"name":        p.Name,
				"description": p.Description,
			}).
			Where(goqu.Ex{"id": p.ID}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, sr.db, query)
}

func (sr *SQLRBAC) FindPermission(ctx context.Context, id uint64) (*domain.Permission, error) {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FindPermission")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Select("id", "name", "description").
			From("permissions").
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Permission](ctx, sr.db, query)
}

func (sr *SQLRBAC) FindPermissionByName(ctx context.Context, name string) (*domain.Permission, error) {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FindPermissionByName")
	defer span.End()

	query := func() (string, []any, error) {
		return sr.qu.Select("id", "name", "description").
			From("permissions").
			Where(goqu.Ex{"name": name}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Permission](ctx, sr.db, query)
}

//nolint:dupl // #12 this is not duplicate with #11
func (sr *SQLRBAC) FetchPermission(ctx context.Context, filter map[string]any) ([]domain.Permission, error) {
	ctx, span := sr.telemetry.Tracer().Start(ctx, "rbac.outbound.SQLRBAC.FetchPermission")
	defer span.End()

	cursor, hasCursor := filter["cursor"].(uint64)
	limit, hasLimit := filter["limit"].(int)
	name, hasName := filter["name"].(string)

	query := func() (string, []any, error) {
		q := sr.qu.Select("id", "name", "description").
			From("permissions")

		if hasCursor && cursor > 0 {
			q = q.Where(goqu.Ex{"id": goqu.Op{"gt": cursor}})
		}

		if hasName {
			q = q.Where(goqu.Ex{"name": goqu.Op{"like": "%" + name + "%"}})
		}

		if hasLimit {
			q = q.Limit(uint(limit + 1))
		}

		return q.Prepared(true).ToSQL()
	}

	return dbops.SQLGets[domain.Permission](ctx, sr.db, query)
}
