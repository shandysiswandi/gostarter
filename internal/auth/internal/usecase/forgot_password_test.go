package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestNewForgotPassword(t *testing.T) {
	type args struct {
		dep Dependency
		s   ForgotPasswordStore
	}
	tests := []struct {
		name string
		args args
		want *ForgotPassword
	}{
		{
			name: "Success",
			args: args{},
			want: &ForgotPassword{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewForgotPassword(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestForgotPassword_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.ForgotPasswordInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ForgotPasswordOutput
		wantErr error
		mockFn  func(a args) *ForgotPassword
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &ForgotPassword{
					telemetry: tel,
					validator: validatorMock,
					idnum:     nil,
					secHash:   nil,
					store:     nil,
					clock:     nil,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmail",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(nil, assert.AnError)

				return &ForgotPassword{
					telemetry: tel,
					validator: validatorMock,
					idnum:     nil,
					secHash:   nil,
					store:     storeMock,
					clock:     nil,
				}
			},
		},
		{
			name: "SuccessButNotFound",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want: &domain.ForgotPasswordOutput{
				Email:   "email",
				Message: msgSuccess,
			},
			wantErr: nil,
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     nil,
					secHash:   nil,
					store:     storeMock,
					clock:     nil,
				}
			},
		},
		{
			name: "ErrorStorePasswordResetByUserID",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       1,
					Name:     "",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				storeMock.EXPECT().
					PasswordResetByUserID(ctx, user.ID).
					Return(nil, assert.AnError)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     nil,
					secHash:   nil,
					store:     storeMock,
					clock:     nil,
				}
			},
		},
		{
			name: "SuccessPasswordResetAlreadyGeneratedButNotExpired",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want: &domain.ForgotPasswordOutput{
				Email:   "email",
				Message: msgSuccess,
			},
			wantErr: nil,
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)
				clockMock := mocker.NewMockClocker(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       1,
					Name:     "",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    user.ID,
					Token:     "token",
					ExpiresAt: now.Add(time.Minute),
				}
				storeMock.EXPECT().
					PasswordResetByUserID(ctx, user.ID).
					Return(ps, nil)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     nil,
					secHash:   nil,
					store:     storeMock,
					clock:     clockMock,
				}
			},
		},
		{
			name: "ErrorStorePasswordResetDelete",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)
				clockMock := mocker.NewMockClocker(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       1,
					Name:     "",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    user.ID,
					Token:     "token",
					ExpiresAt: now.Add(-time.Minute),
				}
				storeMock.EXPECT().
					PasswordResetByUserID(ctx, user.ID).
					Return(ps, nil)

				storeMock.EXPECT().
					PasswordResetDelete(ctx, ps.ID).
					Return(assert.AnError)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     nil,
					secHash:   nil,
					store:     storeMock,
					clock:     clockMock,
				}
			},
		},
		{
			name: "ErrorSecHash",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       1,
					Name:     "",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    user.ID,
					Token:     "token",
					ExpiresAt: now.Add(-time.Minute),
				}
				storeMock.EXPECT().
					PasswordResetByUserID(ctx, user.ID).
					Return(ps, nil)

				storeMock.EXPECT().
					PasswordResetDelete(ctx, ps.ID).
					Return(nil)

				secHashMock.EXPECT().
					Hash(fmt.Sprintf("%d-%v", user.ID, now.Unix())).
					Return(nil, assert.AnError)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     nil,
					secHash:   secHashMock,
					store:     storeMock,
					clock:     clockMock,
				}
			},
		},
		{
			name: "ErrorStorePasswordResetSave",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				idnumMock := mocker.NewMockNumberID(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       1,
					Name:     "",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    user.ID,
					Token:     "token",
					ExpiresAt: now.Add(-time.Minute),
				}
				storeMock.EXPECT().
					PasswordResetByUserID(ctx, user.ID).
					Return(ps, nil)

				storeMock.EXPECT().
					PasswordResetDelete(ctx, ps.ID).
					Return(nil)

				sechashResult := []byte{}
				secHashMock.EXPECT().
					Hash(fmt.Sprintf("%d-%v", user.ID, now.Unix())).
					Return(sechashResult, nil)

				idnumMock.EXPECT().
					Generate().
					Return(111)

				psData := domain.PasswordReset{
					ID:        111,
					UserID:    user.ID,
					Token:     string(sechashResult),
					ExpiresAt: now.Add(time.Hour),
				}
				storeMock.EXPECT().
					PasswordResetSave(ctx, psData).
					Return(assert.AnError)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     idnumMock,
					secHash:   secHashMock,
					store:     storeMock,
					clock:     clockMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  domain.ForgotPasswordInput{Email: "email"},
			},
			want: &domain.ForgotPasswordOutput{
				Email:   "email",
				Message: msgSuccess,
			},
			wantErr: nil,
			mockFn: func(a args) *ForgotPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockForgotPasswordStore(t)
				secHashMock := mocker.NewMockHash(t)
				clockMock := mocker.NewMockClocker(t)
				idnumMock := mocker.NewMockNumberID(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ForgotPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       1,
					Name:     "",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				now := time.Time{}
				clockMock.EXPECT().
					Now().
					Return(now)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    user.ID,
					Token:     "token",
					ExpiresAt: now.Add(-time.Minute),
				}
				storeMock.EXPECT().
					PasswordResetByUserID(ctx, user.ID).
					Return(ps, nil)

				storeMock.EXPECT().
					PasswordResetDelete(ctx, ps.ID).
					Return(nil)

				sechashResult := []byte{}
				secHashMock.EXPECT().
					Hash(fmt.Sprintf("%d-%v", user.ID, now.Unix())).
					Return(sechashResult, nil)

				idnumMock.EXPECT().
					Generate().
					Return(111)

				psData := domain.PasswordReset{
					ID:        111,
					UserID:    user.ID,
					Token:     string(sechashResult),
					ExpiresAt: now.Add(time.Hour),
				}
				storeMock.EXPECT().
					PasswordResetSave(ctx, psData).
					Return(nil)

				return &ForgotPassword{
					telemetry: telemetry.NewTelemetry(),
					validator: validatorMock,
					idnum:     idnumMock,
					secHash:   secHashMock,
					store:     storeMock,
					clock:     clockMock,
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
