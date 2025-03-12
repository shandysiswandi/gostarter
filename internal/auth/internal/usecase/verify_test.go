package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/auth/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestNewVerify(t *testing.T) {
	type args struct {
		dep Dependency
		s   VerifyStore
	}
	tests := []struct {
		name string
		args args
		want *Verify
	}{
		{
			name: "Success",
			args: args{},
			want: &Verify{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewVerify(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVerify_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.VerifyInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.VerifyOutput
		wantErr error
		mockFn  func(a args) *Verify
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in: domain.VerifyInput{
					Email: "email",
					Code:  "code",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *Verify {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Verify{
					tel:       tel,
					validator: validatorMock,
					secHash:   nil,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreUserByEmail",
			args: args{
				ctx: context.Background(),
				in: domain.VerifyInput{
					Email: "email",
					Code:  "code",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Verify {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockVerifyStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(nil, assert.AnError)

				return &Verify{
					tel:       tel,
					validator: validatorMock,
					secHash:   nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreUserByEmailNotFound",
			args: args{
				ctx: context.Background(),
				in: domain.VerifyInput{
					Email: "email",
					Code:  "code",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *Verify {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockVerifyStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				return &Verify{
					tel:       tel,
					validator: validatorMock,
					secHash:   nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "SuccessAlreadyVerified",
			args: args{
				ctx: context.Background(),
				in: domain.VerifyInput{
					Email: "email",
					Code:  "code",
				},
			},
			want: &domain.VerifyOutput{
				Email:    "email",
				VerifyAt: time.Time{},
			},
			wantErr: nil,
			mockFn: func(a args) *Verify {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockVerifyStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "name",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: true},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				return &Verify{
					tel:       tel,
					validator: validatorMock,
					secHash:   nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreUserVerificationByUserID",
			args: args{
				ctx: context.Background(),
				in: domain.VerifyInput{
					Email: "email",
					Code:  "code",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Verify {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockVerifyStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "name",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: false},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				storeMock.EXPECT().
					UserVerificationByUserID(ctx, user.ID).
					Return(nil, assert.AnError)

				return &Verify{
					tel:       tel,
					validator: validatorMock,
					secHash:   nil,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreUserVerificationByUserIDNotFound",
			args: args{
				ctx: context.Background(),
				in: domain.VerifyInput{
					Email: "email",
					Code:  "code",
				},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("Invalid credentials", goerror.CodeUnauthorized),
			mockFn: func(a args) *Verify {
				tel := telemetry.NewTelemetry()
				validatorMock := mocker.NewMockValidator(t)
				storeMock := mockz.NewMockVerifyStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Register")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				user := &domain.User{
					ID:         10,
					Name:       "name",
					Email:      "email",
					Password:   "password",
					VerifiedAt: sql.Null[time.Time]{Valid: false},
				}
				storeMock.EXPECT().
					UserByEmail(ctx, a.in.Email).
					Return(user, nil)

				storeMock.EXPECT().
					UserVerificationByUserID(ctx, user.ID).
					Return(nil, nil)

				return &Verify{
					tel:       tel,
					validator: validatorMock,
					secHash:   nil,
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
