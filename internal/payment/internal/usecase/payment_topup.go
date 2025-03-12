package usecase

import (
	"context"

	"github.com/shandysiswandi/goreng/clock"
	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/goreng/telemetry/logger"
	"github.com/shandysiswandi/goreng/uid"
	"github.com/shandysiswandi/goreng/validation"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/sqlkit"
)

type PaymentTopupStore interface {
	FindAccountByUserID(ctx context.Context, userID uint64) (*domain.Account, error)
	FindTopupByReferenceID(ctx context.Context, refID string) (*domain.Topup, error)
	SaveTopup(ctx context.Context, topup domain.Topup) error
	SaveTransaction(ctx context.Context, topup domain.Transaction) error
	UpdateAccount(ctx context.Context, data map[string]any) error
}

type PaymentTopup struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	clock     clock.Clocker
	trx       dbops.Tx
	store     PaymentTopupStore
}

func NewPaymentTopup(dep Dependency, s PaymentTopupStore) *PaymentTopup {
	return &PaymentTopup{
		telemetry: dep.Telemetry,
		uidnumber: dep.UIDNumber,
		validator: dep.Validator,
		clock:     dep.Clock,
		trx:       dep.Transaction,
		store:     s,
	}
}

func (pt *PaymentTopup) Call(ctx context.Context, in domain.PaymentTopupInput) (
	*domain.PaymentTopupOutput, error,
) {
	ctx, span := pt.telemetry.Tracer().Start(ctx, "payment.usecase.PaymentTopup")
	defer span.End()

	if err := pt.validator.Validate(in); err != nil {
		pt.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("Invalid request payload", err)
	}

	top, err := pt.store.FindTopupByReferenceID(ctx, in.ReferenceID)
	if err != nil {
		pt.telemetry.Logger().Error(ctx, "failed to get topup by ref_id", err,
			logger.KeyVal("reference_id", in.ReferenceID))

		return nil, goerror.NewServerInternal(err)
	}

	if top != nil {
		pt.telemetry.Logger().Warn(ctx, "duplicate request topup by ref_id",
			logger.KeyVal("reference_id", in.ReferenceID))

		return nil, goerror.NewBusiness("duplicate request topup", goerror.CodeConflict)
	}

	clm := lib.GetJWTClaim(ctx)
	acc, err := pt.store.FindAccountByUserID(ctx, clm.AuthID)
	if err != nil {
		pt.telemetry.Logger().Error(ctx, "failed to find account", err, logger.KeyVal("user_id", clm.AuthID))

		return nil, goerror.NewServerInternal(err)
	}

	if acc == nil {
		pt.telemetry.Logger().Warn(ctx, "account is not found", logger.KeyVal("user_id", clm.AuthID))

		return nil, goerror.NewBusiness("account not found", goerror.CodeNotFound)
	}

	if err := pt.doTransaction(ctx, in, acc, clm.AuthID); err != nil {
		return nil, err
	}

	return &domain.PaymentTopupOutput{
		ReferenceID: in.ReferenceID,
		Amount:      in.Amount,
		Balance:     acc.Balanace.Add(in.Amount),
	}, nil
}

func (pt *PaymentTopup) doTransaction(ctx context.Context, in domain.PaymentTopupInput,
	acc *domain.Account, userID uint64,
) error {
	return pt.trx.Transaction(ctx, func(cc context.Context) error {
		trx := domain.Transaction{
			ID:       pt.uidnumber.Generate(),
			UserID:   userID,
			Amount:   in.Amount,
			Type:     domain.TransactionTypeDebit,
			Status:   domain.TransactionStatusPending,
			Remark:   "top up balance",
			CreateAt: pt.clock.Now(),
		}
		if err := pt.store.SaveTransaction(cc, trx); err != nil {
			pt.telemetry.Logger().Error(ctx, "failed to save transaction", err,
				logger.KeyVal("transaction_data", trx))

			return goerror.NewServerInternal(err)
		}

		topup := domain.Topup{
			ID:            pt.uidnumber.Generate(),
			TransactionID: trx.ID,
			ReferenceID:   in.ReferenceID,
			Amount:        in.Amount,
		}
		if err := pt.store.SaveTopup(cc, topup); err != nil {
			pt.telemetry.Logger().Error(ctx, "failed to save topup", err,
				logger.KeyVal("topup_data", topup))

			return goerror.NewServerInternal(err)
		}

		accUpdateData := map[string]any{
			"id":      acc.ID,
			"balance": acc.Balanace.Add(in.Amount),
		}
		if err := pt.store.UpdateAccount(cc, accUpdateData); err != nil {
			pt.telemetry.Logger().Error(ctx, "failed to update account", err,
				logger.KeyVal("account_update_data", accUpdateData))

			return goerror.NewServerInternal(err)
		}

		return nil
	})
}
