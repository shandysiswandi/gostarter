package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/gostarter/internal/payment/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/payment/internal/mockz"
	mclk "github.com/shandysiswandi/gostarter/pkg/clock/mocker"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mu "github.com/shandysiswandi/gostarter/pkg/uid/mocker"
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
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
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
			name: "ErrorStoreFindTopupByReferenceID",
			args: args{
				ctx: context.Background(),
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
			name: "ErrorStoreFindTopupByReferenceIDExists",
			args: args{
				ctx: context.Background(),
				in: domain.PaymentTopupInput{
					ReferenceID: "uuid",
					Amount:      decimal.NewFromFloat(123.45),
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("duplicate request topup", goerror.CodeConflict),
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				out := &domain.Topup{
					ID:            10,
					TransactionID: 20,
					ReferenceID:   "uuid",
					Amount:        decimal.NewFromInt(1000),
				}
				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(out, nil)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
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
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, nil)

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
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, nil)

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
			name: "ErrorTransactionStoreSaveTransaction",
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
				muid := mu.NewMockNumberID(t)
				clk := mclk.NewMockClocker(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, nil)

				account := &domain.Account{
					ID:       22,
					UserID:   11,
					Balanace: decimal.NewFromInt(1000),
				}
				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(account, nil)

				ctx = dbops.SetContextNoopTx(ctx)

				muid.EXPECT().
					Generate().
					Return(16)

				clk.EXPECT().
					Now().
					Return(time.Time{})

				dataTrx := domain.Transaction{
					ID:       16,
					UserID:   11,
					Amount:   decimal.NewFromFloat(123.45),
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "top up balance",
					CreateAt: time.Time{},
				}
				storeMock.EXPECT().
					SaveTransaction(ctx, dataTrx).
					Return(assert.AnError)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
					trx:       trxMock,
					uidnumber: muid,
					clock:     clk,
				}
			},
		},
		{
			name: "ErrorTransactionStoreSaveTopup",
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
				muid := mu.NewMockNumberID(t)
				clk := mclk.NewMockClocker(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, nil)

				account := &domain.Account{
					ID:       22,
					UserID:   11,
					Balanace: decimal.NewFromInt(1000),
				}
				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(account, nil)

				ctx = dbops.SetContextNoopTx(ctx)

				muid.EXPECT().
					Generate().
					Return(16).
					Once()

				clk.EXPECT().
					Now().
					Return(time.Time{})

				dataTrx := domain.Transaction{
					ID:       16,
					UserID:   11,
					Amount:   a.in.Amount,
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "top up balance",
					CreateAt: time.Time{},
				}
				storeMock.EXPECT().
					SaveTransaction(ctx, dataTrx).
					Return(nil)

				muid.EXPECT().
					Generate().
					Return(19).
					Once()

				dataTopup := domain.Topup{
					ID:            19,
					TransactionID: dataTrx.ID,
					ReferenceID:   a.in.ReferenceID,
					Amount:        a.in.Amount,
				}
				storeMock.EXPECT().
					SaveTopup(ctx, dataTopup).
					Return(assert.AnError)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
					trx:       trxMock,
					uidnumber: muid,
					clock:     clk,
				}
			},
		},
		{
			name: "ErrorTransactionStoreUpdateAccount",
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
				muid := mu.NewMockNumberID(t)
				clk := mclk.NewMockClocker(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, nil)

				account := &domain.Account{
					ID:       22,
					UserID:   11,
					Balanace: decimal.NewFromInt(1000),
				}
				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(account, nil)

				ctx = dbops.SetContextNoopTx(ctx)

				muid.EXPECT().
					Generate().
					Return(16).
					Once()

				clk.EXPECT().
					Now().
					Return(time.Time{})

				dataTrx := domain.Transaction{
					ID:       16,
					UserID:   11,
					Amount:   a.in.Amount,
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "top up balance",
					CreateAt: time.Time{},
				}
				storeMock.EXPECT().
					SaveTransaction(ctx, dataTrx).
					Return(nil).
					Once()

				muid.EXPECT().
					Generate().
					Return(19)

				dataTopup := domain.Topup{
					ID:            19,
					TransactionID: dataTrx.ID,
					ReferenceID:   a.in.ReferenceID,
					Amount:        a.in.Amount,
				}
				storeMock.EXPECT().
					SaveTopup(ctx, dataTopup).
					Return(nil)

				dataUpdateAccount := map[string]any{
					"id":      account.ID,
					"balance": account.Balanace.Add(a.in.Amount),
				}
				storeMock.EXPECT().
					UpdateAccount(ctx, dataUpdateAccount).
					Return(assert.AnError)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
					trx:       trxMock,
					uidnumber: muid,
					clock:     clk,
				}
			},
		},
		{
			name: "Success",
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
				Balance:     decimal.NewFromFloat(123.45).Add(decimal.NewFromInt(1000)),
			},
			wantErr: nil,
			mockFn: func(a args) *PaymentTopup {
				tel := telemetry.NewTelemetry()
				validatorMock := mv.NewMockValidator(t)
				muid := mu.NewMockNumberID(t)
				clk := mclk.NewMockClocker(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())
				storeMock := mockz.NewMockPaymentTopupStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "payment.usecase.PaymentTopup")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindTopupByReferenceID(ctx, a.in.ReferenceID).
					Return(nil, nil)

				account := &domain.Account{
					ID:       22,
					UserID:   11,
					Balanace: decimal.NewFromInt(1000),
				}
				storeMock.EXPECT().
					FindAccountByUserID(ctx, uint64(11)).
					Return(account, nil)

				ctx = dbops.SetContextNoopTx(ctx)

				muid.EXPECT().
					Generate().
					Return(16).
					Once()

				clk.EXPECT().
					Now().
					Return(time.Time{})

				dataTrx := domain.Transaction{
					ID:       16,
					UserID:   11,
					Amount:   a.in.Amount,
					Type:     domain.TransactionTypeDebit,
					Status:   domain.TransactionStatusPending,
					Remark:   "top up balance",
					CreateAt: time.Time{},
				}
				storeMock.EXPECT().
					SaveTransaction(ctx, dataTrx).
					Return(nil).
					Once()

				muid.EXPECT().
					Generate().
					Return(19)

				dataTopup := domain.Topup{
					ID:            19,
					TransactionID: dataTrx.ID,
					ReferenceID:   a.in.ReferenceID,
					Amount:        a.in.Amount,
				}
				storeMock.EXPECT().
					SaveTopup(ctx, dataTopup).
					Return(nil)

				dataUpdateAccount := map[string]any{
					"id":      account.ID,
					"balance": account.Balanace.Add(a.in.Amount),
				}
				storeMock.EXPECT().
					UpdateAccount(ctx, dataUpdateAccount).
					Return(nil)

				return &PaymentTopup{
					telemetry: tel,
					validator: validatorMock,
					store:     storeMock,
					trx:       trxMock,
					uidnumber: muid,
					clock:     clk,
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
