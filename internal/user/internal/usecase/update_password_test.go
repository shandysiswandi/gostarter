package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	mockHash "github.com/shandysiswandi/goreng/mocker"
	mockValidation "github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdatePassword(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		s    UpdatePasswordStore
		want *UpdatePassword
	}{
		{
			name: "Success",
			dep:  Dependency{},
			s:    nil,
			want: &UpdatePassword{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUpdatePassword(tt.dep, tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdatePassword_Call(t *testing.T) {
	claim := lib.NewJWTClaim(11, "email", time.Time{}, nil)
	ctxJWT := lib.SetJWTClaim(context.Background(), claim)

	type args struct {
		ctx context.Context
		in  domain.UpdatePasswordInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr error
		mockFn  func(a args) *UpdatePassword
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &UpdatePassword{
					tel:       tel,
					validator: validatorMock,
					hash:      nil,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreFindUser",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockUpdatePasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUser(ctx, uint64(11)).
					Return(nil, assert.AnError)

				return &UpdatePassword{
					tel:       tel,
					hash:      nil,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindUserNotFound",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockUpdatePasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUser(ctx, uint64(11)).
					Return(nil, nil)

				return &UpdatePassword{
					tel:       tel,
					hash:      nil,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorComparePassword",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				hashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockUpdatePasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       11,
					Name:     "test",
					Email:    "test@test.com",
					Password: "***",
				}
				storeMock.EXPECT().
					FindUser(ctx, uint64(11)).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.CurrentPassword).
					Return(false)

				return &UpdatePassword{
					tel:       tel,
					hash:      hashMock,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorHashPassword",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				hashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockUpdatePasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       11,
					Name:     "test",
					Email:    "test@test.com",
					Password: "***",
				}
				storeMock.EXPECT().
					FindUser(ctx, uint64(11)).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.CurrentPassword).
					Return(true)

				hashMock.EXPECT().
					Hash(a.in.NewPassword).
					Return(nil, assert.AnError)

				return &UpdatePassword{
					tel:       tel,
					hash:      hashMock,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreUpdatePassword",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				hashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockUpdatePasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       11,
					Name:     "test",
					Email:    "test@test.com",
					Password: "***",
				}
				storeMock.EXPECT().
					FindUser(ctx, uint64(11)).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.CurrentPassword).
					Return(true)

				hashMock.EXPECT().
					Hash(a.in.NewPassword).
					Return([]byte("hash"), nil)

				user1 := *user
				user1.Password = "hash"
				storeMock.EXPECT().
					UpdatePassword(ctx, user1).
					Return(assert.AnError)

				return &UpdatePassword{
					tel:       tel,
					hash:      hashMock,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: ctxJWT,
				in: domain.UpdatePasswordInput{
					CurrentPassword: "old password",
					NewPassword:     "new password",
				},
			},
			want: &domain.User{
				ID:       11,
				Name:     "test",
				Email:    "test@test.com",
				Password: "***",
			},
			wantErr: nil,
			mockFn: func(a args) *UpdatePassword {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				hashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockUpdatePasswordStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.UpdatePassword")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:       11,
					Name:     "test",
					Email:    "test@test.com",
					Password: "***",
				}
				storeMock.EXPECT().
					FindUser(ctx, uint64(11)).
					Return(user, nil)

				hashMock.EXPECT().
					Verify(user.Password, a.in.CurrentPassword).
					Return(true)

				hashMock.EXPECT().
					Hash(a.in.NewPassword).
					Return([]byte("hash"), nil)

				user1 := *user
				user1.Password = "hash"
				storeMock.EXPECT().
					UpdatePassword(ctx, user1).
					Return(nil)

				return &UpdatePassword{
					tel:       tel,
					hash:      hashMock,
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
