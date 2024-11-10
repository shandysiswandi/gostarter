package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type SQLAuth struct {
	db        *sql.DB
	qu        goqu.DialectWrapper
	telemetry *telemetry.Telemetry
}

func NewSQLAuth(db *sql.DB, qu goqu.DialectWrapper, tel *telemetry.Telemetry) *SQLAuth {
	return &SQLAuth{
		db:        db,
		qu:        qu,
		telemetry: tel,
	}
}

func (st *SQLAuth) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "FindUserByEmail")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "email", "password").
			From("users").
			Where(goqu.Ex{"email": email}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.User](ctx, st.db, query)
}

func (st *SQLAuth) SaveUser(ctx context.Context, u domain.User) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "SaveUser")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Insert("users").
			Cols("id", "email", "password").
			Vals([]any{u.ID, u.Email, u.Password}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrUserNotCreated
	}

	return err
}

func (st *SQLAuth) UpdateUserPassword(ctx context.Context, id uint64, pass string) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "UpdateUserPassword")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Update("users").
			Set(map[string]any{"password": pass}).
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}

func (st *SQLAuth) FindTokenByUserID(ctx context.Context, uid uint64) (*domain.Token, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "FindTokenByUserID")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select(
			"id",
			"user_id",
			"access_token",
			"refresh_token",
			"access_expires_at",
			"refresh_expires_at",
		).
			From("tokens").
			Where(goqu.Ex{"user_id": uid}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Token](ctx, st.db, query)
}

func (st *SQLAuth) FindTokenByRefresh(ctx context.Context, ref string) (*domain.Token, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "FindTokenByRefresh")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select(
			"id",
			"user_id",
			"access_token",
			"refresh_token",
			"access_expires_at",
			"refresh_expires_at",
		).
			From("tokens").
			Where(goqu.Ex{"refresh_token": ref}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Token](ctx, st.db, query)
}

func (st *SQLAuth) SaveToken(ctx context.Context, t domain.Token) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "SaveToken")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Insert("tokens").
			Cols(
				"id",
				"user_id",
				"access_token",
				"refresh_token",
				"access_expires_at",
				"refresh_expires_at",
			).
			Vals([]any{
				t.ID,
				t.UserID,
				t.AccessToken,
				t.RefreshToken,
				t.AccessExpiredAt,
				t.RefreshExpiredAt,
			}).
			OnConflict(
				goqu.DoUpdate("user_id", goqu.Record{
					"access_token":       t.AccessToken,
					"refresh_token":      t.RefreshToken,
					"access_expires_at":  t.AccessExpiredAt,
					"refresh_expires_at": t.RefreshExpiredAt,
				}),
			).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrTokenNoRowsAffected
	}

	return err
}

func (st *SQLAuth) FindPasswordResetByUserID(ctx context.Context, uid uint64) (*domain.PasswordReset, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "FindPasswordResetByUserID")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "user_id", "token", "expires_at").
			From("password_resets").
			Where(goqu.Ex{"user_id": uid}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.PasswordReset](ctx, st.db, query)
}

func (st *SQLAuth) FindPasswordResetByToken(ctx context.Context, t string) (*domain.PasswordReset, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "FindPasswordResetByToken")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "user_id", "token", "expires_at").
			From("password_resets").
			Where(goqu.Ex{"token": t}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.PasswordReset](ctx, st.db, query)
}

func (st *SQLAuth) SavePasswordReset(ctx context.Context, ps domain.PasswordReset) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "SavePasswordReset")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Insert("password_resets").
			Cols("id", "user_id", "token", "expires_at").
			Vals([]any{ps.ID, ps.UserID, ps.Token, ps.ExpiresAt}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrPasswordResetNotCreated
	}

	return err
}

func (st *SQLAuth) DeletePasswordReset(ctx context.Context, id uint64) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "DeletePasswordReset")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Delete("password_resets").
			Where(goqu.Ex{"id": id}).
			Prepared(true).
			ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}
