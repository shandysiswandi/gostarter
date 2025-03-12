package outbound

import (
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
	"github.com/shopspring/decimal"
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
	ctx, span := st.telemetry.Tracer().Start(ctx, "payment.outbound.SQLPayment.FindAccountByUserID")
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
	ctx, span := st.telemetry.Tracer().Start(ctx, "payment.outbound.SQLPayment.FindTopupByReferenceID")
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

func (st *SQLPayment) SaveTopup(ctx context.Context, t domain.Topup) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "payment.outbound.SQLPayment.SaveTopup")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Insert("topups").
			Cols("id", "transaction_id", "reference_id", "amount").
			Vals([]any{t.ID, t.TransactionID, t.ReferenceID, t.Amount}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrAccountNoRowsAffected
	}

	return err
}

func (st *SQLPayment) SaveTransaction(ctx context.Context, t domain.Transaction) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "payment.outbound.SQLPayment.SaveTransaction")
	defer span.End()

	query := func() (string, []any, error) {
		return st.qu.Insert("transactions").
			Cols("id", "user_id", "amount", "type", "status", "remark", "created_at").
			Vals([]any{t.ID, t.UserID, t.Amount, t.Type, t.Status, t.Remark, t.CreateAt}).
			Prepared(true).
			ToSQL()
	}

	err := dbops.Exec(ctx, st.db, query, true)
	if errors.Is(err, dbops.ErrZeroRowsAffected) {
		return domain.ErrTransactionNoRowsAffected
	}

	return err
}

func (st *SQLPayment) UpdateAccount(ctx context.Context, data map[string]any) error {
	ctx, span := st.telemetry.Tracer().Start(ctx, "payment.outbound.SQLPayment.UpdateAccount")
	defer span.End()

	id, hasID := data["id"].(uint64)
	balance, hasBalance := data["balance"].(decimal.Decimal)
	sets := map[string]any{}

	query := func() (string, []any, error) {
		q := st.qu.Update("accounts")

		if hasBalance {
			sets["balance"] = balance
		}

		if len(sets) > 0 {
			q = q.Set(sets)
		}

		if hasID {
			q = q.Where(goqu.Ex{"id": id})
		}

		return q.Prepared(true).ToSQL()
	}

	return dbops.Exec(ctx, st.db, query)
}
