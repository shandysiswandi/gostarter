package usecase

import (
	"context"

	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/telemetry/logger"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

type PaymentTopupStore interface {
	FindAccountByUserID(ctx context.Context, userID uint64) (*domain.Account, error)
	FindTopupByReferenceID(ctx context.Context, refID string) (*domain.Topup, error)
}

type PaymentTopup struct {
	telemetry *telemetry.Telemetry
	validator validation.Validator
	uidnumber uid.NumberID
	trx       dbops.Tx
	store     PaymentTopupStore
}

func NewPaymentTopup(dep Dependency, s PaymentTopupStore) *PaymentTopup {
	return &PaymentTopup{
		telemetry: dep.Telemetry,
		uidnumber: dep.UIDNumber,
		validator: dep.Validator,
		trx:       dep.Transaction,
		store:     s,
	}
}

func (pt *PaymentTopup) Call(ctx context.Context, in domain.PaymentTopupInput) (
	*domain.PaymentTopupOutput, error,
) {
	if err := pt.validator.Validate(in); err != nil {
		pt.telemetry.Logger().Warn(ctx, "validation failed")

		return nil, goerror.NewInvalidInput("validation input fail", err)
	}

	clm := jwt.GetClaim(ctx)

	acc, err := pt.store.FindAccountByUserID(ctx, clm.AuthID)
	if err != nil {
		pt.telemetry.Logger().Error(ctx, "account fail to find", err, logger.KeyVal("user_id", clm.AuthID))

		return nil, goerror.NewServer("failed to find account", err)
	}

	if acc == nil {
		pt.telemetry.Logger().Warn(ctx, "account is not found", logger.KeyVal("user_id", clm.AuthID))

		return nil, goerror.NewBusiness("account not found", goerror.CodeNotFound)
	}

	top, err := pt.store.FindTopupByReferenceID(ctx, in.ReferenceID)
	if err != nil {
		pt.telemetry.Logger().Error(ctx, "todo fail to create", err)

		return nil, goerror.NewServer("failed to create todo", err)
	}

	if top != nil {
		return &domain.PaymentTopupOutput{
			ReferenceID: in.ReferenceID,
			Amount:      in.Amount,
			Balance:     acc.Balanace,
		}, nil
	}

	err = pt.trx.Transaction(ctx, func(_ context.Context) error {
		// update balance account
		// create topUps
		// create transaction
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &domain.PaymentTopupOutput{
		ReferenceID: in.ReferenceID,
		Amount:      in.Amount,
		Balance:     acc.Balanace.Add(in.Amount),
	}, nil
}
