package outbound

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type SQLPayment struct {
	db        *sql.DB
	qu        goqu.DialectWrapper
	telemetry *telemetry.Telemetry
}

func NewSQLPayment(db *sql.DB, qu goqu.DialectWrapper, tel *telemetry.Telemetry) *SQLPayment {
	return &SQLPayment{
		db:        db,
		qu:        qu,
		telemetry: tel,
	}
}

func (st *SQLPayment) FindAccountByUserID(ctx context.Context, userID uint64) (*domain.Account, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.FindAccountByUserID")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "user_id", "balance").
			From("accounts").
			Where(goqu.Ex{"user_id": userID}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Account](ctx, st.db, query)
}

func (st *SQLPayment) FindTopupByReferenceID(ctx context.Context, refID string) (*domain.Topup, error) {
	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.FindTopupByReferenceID")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Select("id", "transaction_id", "reference_id", "amount").
			From("topups").
			Where(goqu.Ex{"reference_id": refID}).
			Prepared(true).
			ToSQL()
	}

	return dbops.SQLGet[domain.Topup](ctx, st.db, query)
}

// func (st *SQLPayment) SaveToken(ctx context.Context, t domain.Token) error {
// 	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.SaveToken")
// 	defer span.End()

// 	query := func() (string, []any, error) {
// 		return st.qu.Insert("tokens").
// 			Cols(
// 				"id",
// 				"user_id",
// 				"access_token",
// 				"refresh_token",
// 				"access_expires_at",
// 				"refresh_expires_at",
// 			).
// 			Vals([]any{
// 				t.ID,
// 				t.UserID,
// 				t.AccessToken,
// 				t.RefreshToken,
// 				t.AccessExpiredAt,
// 				t.RefreshExpiredAt,
// 			}).
// 			Prepared(true).
// 			ToSQL()
// 	}

// 	err := dbops.Exec(ctx, st.db, query, true)
// 	if errors.Is(err, dbops.ErrZeroRowsAffected) {
// 		return domain.ErrTokenNoRowsAffected
// 	}

// 	return err
// }

// func (st *SQLPayment) UpdateToken(ctx context.Context, t domain.Token) error {
// 	ctx, span := st.telemetry.Tracer().Start(ctx, "outbound.UpdateToken")
// 	defer span.End()

// 	query := func() (string, []any, error) {
// 		return st.qu.Update("tokens").
// 			Set(map[string]any{
// 				"user_id":            t.UserID,
// 				"access_token":       t.AccessToken,
// 				"refresh_token":      t.RefreshToken,
// 				"access_expires_at":  t.AccessExpiredAt,
// 				"refresh_expires_at": t.RefreshExpiredAt,
// 			}).
// 			Where(goqu.Ex{"id": t.ID}).
// 			Prepared(true).
// 			ToSQL()
// 	}

// 	return dbops.Exec(ctx, st.db, query)
// }
