package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/dbops"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	mockHash "github.com/shandysiswandi/gostarter/pkg/hash/mocker"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mockUID "github.com/shandysiswandi/gostarter/pkg/uid/mocker"
	mockValidation "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewRegister(t *testing.T) {
	type args struct {
		dep Dependency
		s   RegisterStore
	}
	tests := []struct {
		name string
		args args
		want *Register
	}{
		{
			name: "Success",
			args: args{},
			want: &Register{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewRegister(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegister_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.RegisterInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.RegisterOutput
		wantErr error
		mockFn  func(a args) *Register
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Name:     "name",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: nil,
					hash:      nil,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmail",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Name:     "name",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockRegisterStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, assert.AnError)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: nil,
					hash:      nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmailAlreadyExists",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Name:     "name",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("email already registered", goerror.CodeConflict),
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockRegisterStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       10,
					Name:     "name",
					Email:    a.in.Email,
					Password: "***",
				}
				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(user, nil)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: nil,
					hash:      nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorHash",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Name:     "name",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockRegisterStore(t)
				hashMock := mockHash.NewMockHash(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return(nil, assert.AnError)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: nil,
					hash:      hashMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorTransactionStoreSaveUser",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Name:     "name",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockRegisterStore(t)
				hashMock := mockHash.NewMockHash(t)
				idnumMock := mockUID.NewMockNumberID(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return([]byte("hash_password"), nil)

				idnumMock.EXPECT().
					Generate().
					Return(111)

				ctx = dbops.SetContextNoopTx(ctx)

				dataUser := domain.User{
					ID:       111,
					Name:     a.in.Name,
					Email:    a.in.Email,
					Password: "hash_password",
				}
				storeMock.EXPECT().
					SaveUser(ctx, dataUser).
					Return(assert.AnError)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: idnumMock,
					hash:      hashMock,
					trx:       trxMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorTransactionStoreSaveAccount",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Name:     "name",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockRegisterStore(t)
				hashMock := mockHash.NewMockHash(t)
				idnumMock := mockUID.NewMockNumberID(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return([]byte("hash_password"), nil)

				idnumMock.EXPECT().
					Generate().
					Return(111).
					Once()

				ctx = dbops.SetContextNoopTx(ctx)

				dataUser := domain.User{
					ID:       111,
					Name:     a.in.Name,
					Email:    a.in.Email,
					Password: "hash_password",
				}
				storeMock.EXPECT().
					SaveUser(ctx, dataUser).
					Return(nil)

				idnumMock.EXPECT().
					Generate().
					Return(121).
					Once()

				dataAccount := domain.Account{
					ID:     121,
					UserID: dataUser.ID,
				}
				storeMock.EXPECT().
					SaveAccount(ctx, dataAccount).
					Return(assert.AnError)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: idnumMock,
					hash:      hashMock,
					trx:       trxMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: domain.RegisterInput{
					Email:    "email",
					Password: "password",
				},
			},
			want:    &domain.RegisterOutput{Email: "email"},
			wantErr: nil,
			mockFn: func(a args) *Register {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockRegisterStore(t)
				hashMock := mockHash.NewMockHash(t)
				idnumMock := mockUID.NewMockNumberID(t)
				trxMock := dbops.NewTransaction(dbops.NewNoopDB())

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return([]byte("hash_password"), nil)

				idnumMock.EXPECT().
					Generate().
					Return(111).
					Once()

				ctx = dbops.SetContextNoopTx(ctx)

				dataUser := domain.User{
					ID:       111,
					Name:     a.in.Name,
					Email:    a.in.Email,
					Password: "hash_password",
				}
				storeMock.EXPECT().
					SaveUser(ctx, dataUser).
					Return(nil)

				idnumMock.EXPECT().
					Generate().
					Return(121).
					Once()

				dataAccount := domain.Account{
					ID:     121,
					UserID: dataUser.ID,
				}
				storeMock.EXPECT().
					SaveAccount(ctx, dataAccount).
					Return(nil)

				return &Register{
					tele:      tel,
					validator: validatorMock,
					uidnumber: idnumMock,
					hash:      hashMock,
					trx:       trxMock,
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
