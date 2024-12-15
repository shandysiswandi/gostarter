package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	mockHash "github.com/shandysiswandi/gostarter/pkg/hash/mocker"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mockValidation "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewResetPassword(t *testing.T) {
	type args struct {
		dep Dependency
		s   ResetPasswordStore
	}
	tests := []struct {
		name string
		args args
		want *ResetPassword
	}{
		{
			name: "Success",
			args: args{},
			want: &ResetPassword{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewResetPassword(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResetPassword_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.ResetPasswordInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ResetPasswordOutput
		wantErr error
		mockFn  func(a args) *ResetPassword
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      nil,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreFindPasswordResetByToken",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(nil, assert.AnError)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindPasswordResetByTokenNotFound",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("invalid token", goerror.CodeUnauthorized),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(nil, nil)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorPasswordResetIsExpired",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("token has expired", goerror.CodeUnauthorized),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    20,
					Token:     "token",
					ExpiresAt: time.Time{},
				}
				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(ps, nil)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreDeletePasswordReset",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    20,
					Token:     "token",
					ExpiresAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(ps, nil)

				storeMock.EXPECT().
					DeletePasswordReset(ctx, ps.ID).
					Return(assert.AnError)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorHashPassword",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)
				hashMock := mockHash.NewMockHash(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    20,
					Token:     "token",
					ExpiresAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(ps, nil)

				storeMock.EXPECT().
					DeletePasswordReset(ctx, ps.ID).
					Return(nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return(nil, assert.AnError)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      hashMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreUpdateUserPassword",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)
				hashMock := mockHash.NewMockHash(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    20,
					Token:     "token",
					ExpiresAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(ps, nil)

				storeMock.EXPECT().
					DeletePasswordReset(ctx, ps.ID).
					Return(nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return([]byte("hash_password"), nil)

				storeMock.EXPECT().
					UpdateUserPassword(ctx, ps.UserID, "hash_password").
					Return(assert.AnError)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      hashMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: domain.ResetPasswordInput{
					Token:    "token",
					Password: "password",
				},
			},
			want: &domain.ResetPasswordOutput{
				Message: "Your password has been successfully reset.",
			},
			wantErr: nil,
			mockFn: func(a args) *ResetPassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockResetPasswordStore(t)
				hashMock := mockHash.NewMockHash(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.ResetPassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				ps := &domain.PasswordReset{
					ID:        10,
					UserID:    20,
					Token:     "token",
					ExpiresAt: time.Now().Add(time.Minute),
				}
				storeMock.EXPECT().
					FindPasswordResetByToken(ctx, a.in.Token).
					Return(ps, nil)

				storeMock.EXPECT().
					DeletePasswordReset(ctx, ps.ID).
					Return(nil)

				hashMock.EXPECT().
					Hash(a.in.Password).
					Return([]byte("hash_password"), nil)

				storeMock.EXPECT().
					UpdateUserPassword(ctx, ps.UserID, "hash_password").
					Return(nil)

				return &ResetPassword{
					telemetry: tel,
					validator: validatorMock,
					hash:      hashMock,
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
