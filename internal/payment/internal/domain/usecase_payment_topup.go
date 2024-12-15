package domain

import (
	"context"

	"github.com/shopspring/decimal"
)

type PaymentTopup interface {
	Call(ctx context.Context, in PaymentTopupInput) (*PaymentTopupOutput, error)
}

type PaymentTopupInput struct {
	ReferenceID string          `validate:"required"`
	Amount      decimal.Decimal `validate:"required"`
}

type PaymentTopupOutput struct {
	ReferenceID string
	Amount      decimal.Decimal
	Balance     decimal.Decimal
}
