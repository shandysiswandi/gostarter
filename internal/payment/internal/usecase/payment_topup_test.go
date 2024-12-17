package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mv "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewPaymentTopup(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		s    PaymentTopupStore
		want *PaymentTopup
	}{
		{
			name: "Success",
			dep:  Dependency{},
			s:    nil,
			want: &PaymentTopup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewPaymentTopup(tt.dep, tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPaymentTopup_Call(t *testing.T) {
	claim := jwt.NewClaim(11, "email", time.Time{}, nil)
	ctxJWT := jwt.SetClaim(context.Background(), claim)

	type args struct {
		ctx context.Context
		in  domain.PaymentTopupInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.PaymentTopupOutput
		wantErr error
		mockFn  func(a args) *PaymentTopup
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in: domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      decimal.NewFromFloat(123.45),
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreFindAccountByUserID",
			args: args{
				ctx: ctxJWT,
				in: domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      decimal.NewFromFloat(123.45),
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(nil, assert.AnError)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindAccountByUserIDNotFound",
			args: args{
				ctx: ctxJWT,
				in: domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      decimal.NewFromFloat(123.45),
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("account not found", goerror.CodeNotFound),
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(nil, nil)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindTopupByReferenceID",
			args: args{
				ctx: ctxJWT,
				in: domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      decimal.NewFromFloat(123.45),
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				account := &domain.Account{
					ID:       89,
					UserID:   11,
					Balanace: decimal.NewFromInt(900),
				}
				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(account, nil)

				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, assert.AnError)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "SuccessStoreFindTopupByReferenceIDExists",
			args: args{
				ctx: ctxJWT,
				in: domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      decimal.NewFromFloat(123.45),
				},
			},
			want: &domain.PaymentTopupOutput{
				ReferenceID: "uuid",
				Amount:      decimal.NewFromFloat(123.45),
				Balance:     decimal.NewFromFloat(900),
			},
			wantErr: nil,
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				account := &domain.Account{
					ID:       89,
					UserID:   11,
					Balanace: decimal.NewFromFloat(900),
				}
				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(account, nil)

				topup := &domain.Topup{
					ID:            33,
					TransactionID: 123,
					ReferenceID:   "uuid",
					Amount:        decimal.NewFromFloat(123),
				}
				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(topup, nil)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := tt.mockFn(tt.args)
			got, err := s.Call(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
