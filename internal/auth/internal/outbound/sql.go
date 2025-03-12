package outbound

import (
	"context"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type SQL struct {
	db        *sqlkit.DB
	telemetry *telemetry.Telemetry
}

func NewSQL(db *sqlkit.DB, tel *telemetry.Telemetry) *SQL {
	return &SQL{
		db:        db,
		telemetry: tel,
	}
}

/*
 * Table: users
 */

// UserByEmail is sql store for get data from table users
func (s *SQL) UserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.UserByEmail")
	defer span.End()

	return sqlkit.One[domain.User](ctx, s.db, sqlkit.Ex{"email": email})
}

// UserSave is sql store for save data to table users
func (s *SQL) UserSave(ctx context.Context, u domain.User) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.UserSave")
	defer span.End()

	query := `INSERT INTO users(id, name, email, password) VALUES(?,?,?,?);`
	args := []any{u.ID, u.Name, u.Email, u.Password}

	result, err := sqlkit.Exec(ctx, s.db, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotCreated
	}

	return nil
}

// UserUpdatePassword is sql store for update data to table users
func (s *SQL) UserUpdatePassword(ctx context.Context, id uint64, password string) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.UserUpdatePassword")
	defer span.End()

	query := `UPDATE users SET password=? WHERE id=?;`
	args := []any{password, id}

	_, err := sqlkit.Exec(ctx, s.db, query, args...)

	return err
}

/*
 * Table: user_verifications
 */

// UserVerificationByUserID is sql store for get data from table user_verifications
func (s *SQL) UserVerificationByUserID(
	ctx context.Context,
	uid uint64,
) (*domain.UserVerification, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.UserVerificationByUserID")
	defer span.End()

	return sqlkit.One[domain.UserVerification](ctx, s.db, sqlkit.Ex{"user_id": uid})
}

/*
 * Table: accounts
 */

// AccountSave is sql store for save data to table accounts
func (s *SQL) AccountSave(ctx context.Context, a domain.Account) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.AccountSave")
	defer span.End()

	query := `INSERT INTO accounts(id, user_id) VALUES(?,?);`
	args := []any{a.ID, a.UserID}

	result, err := sqlkit.Exec(ctx, s.db, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrAccountNotCreated
	}

	return nil
}

/*
 * Table: tokens
 */

// TokenByUserID is sql store for get data from table tokens
func (s *SQL) TokenByUserID(ctx context.Context, uid uint64) (*domain.Token, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.TokenByUserID")
	defer span.End()

	return sqlkit.One[domain.Token](ctx, s.db, sqlkit.Ex{"user_id": uid})
}

// TokenByRefresh is sql store for get data from table tokens
func (s *SQL) TokenByRefresh(ctx context.Context, refresh string) (*domain.Token, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.TokenByRefresh")
	defer span.End()

	return sqlkit.One[domain.Token](ctx, s.db, sqlkit.Ex{"refresh_token": refresh})
}

// TokenSave is sql store for save data to table tokens
func (s *SQL) TokenSave(ctx context.Context, t domain.Token) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.TokenSave")
	defer span.End()

	query := `INSERT INTO tokens(id, user_id, access_token, refresh_token, access_expires_at,
	refresh_expires_at) VALUES(?, ?, ?, ?, ?, ?);`
	args := []any{
		t.ID,
		t.UserID,
		t.AccessToken,
		t.RefreshToken,
		t.AccessExpiresAt,
		t.RefreshExpiresAt,
	}

	result, err := sqlkit.Exec(ctx, s.db, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrTokenNoRowsAffected
	}

	return nil
}

// TokenUpdate is sql store for update data to table tokens
func (s *SQL) TokenUpdate(ctx context.Context, t domain.Token) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.TokenUpdate")
	defer span.End()

	query := `UPDATE tokens SET user_id=?, access_token=?, refresh_token=?,
	access_expires_at=?, refresh_expires_at=? WHERE id=?;`
	args := []any{t.UserID, t.AccessToken, t.RefreshToken, t.AccessExpiresAt, t.RefreshExpiresAt, t.ID}

	_, err := sqlkit.Exec(ctx, s.db, query, args...)

	return err
}

/*
 * Table: password_resets
 */

// PasswordResetByUserID is sql store for get data from table password_resets
func (s *SQL) PasswordResetByUserID(ctx context.Context, uid uint64) (*domain.PasswordReset, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.PasswordResetByUserID")
	defer span.End()

	return sqlkit.One[domain.PasswordReset](ctx, s.db, sqlkit.Ex{"user_id": uid})
}

// PasswordResetByToken is sql store for get data from table password_resets
func (s *SQL) PasswordResetByToken(ctx context.Context, token string) (*domain.PasswordReset, error) {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.PasswordResetByToken")
	defer span.End()

	return sqlkit.One[domain.PasswordReset](ctx, s.db, sqlkit.Ex{"token": token})

	// return dbops.SQLGet[domain.PasswordReset](ctx, st.db, query)
}

// PasswordResetSave is sql store for save data to table password_resets
func (s *SQL) PasswordResetSave(ctx context.Context, ps domain.PasswordReset) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.PasswordResetSave")
	defer span.End()

	query := `INSERT INTO password_resets(id, user_id, token, expires_at) VALUES(?, ?, ?, ?);`
	args := []any{ps.ID, ps.UserID, ps.Token, ps.ExpiresAt}

	result, err := sqlkit.Exec(ctx, s.db, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return domain.ErrPasswordResetNotCreated
	}

	return nil
}

// PasswordResetDelete is sql store for delete data to table password_resets
func (s *SQL) PasswordResetDelete(ctx context.Context, id uint64) error {
	ctx, span := s.telemetry.Tracer().Start(ctx, "auth.outbound.SQL.PasswordResetDelete")
	defer span.End()

	query := `DELETE password_resets WHERE id=?;`
	args := []any{id}

	_, err := sqlkit.Exec(ctx, s.db, query, args...)

	return err
}
