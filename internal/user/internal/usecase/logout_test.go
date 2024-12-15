package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	mockHash "github.com/shandysiswandi/gostarter/pkg/hash/mocker"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mockValidation "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewLogout(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		s    LogoutStore
		want *Logout
	}{
		{
			name: "Success",
			dep:  Dependency{},
			s:    nil,
			want: &Logout{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewLogout(tt.dep, tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLogout_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.LogoutInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.LogoutOutput
		wantErr error
		mockFn  func(a args) *Logout
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in:  domain.LogoutInput{AccessToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Logout {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.Logout")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Logout{
					tel:       tel,
					validator: validatorMock,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorSecHash",
			args: args{
				ctx: context.Background(),
				in:  domain.LogoutInput{AccessToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Logout {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)

				_, span := tel.Tracer().Start(a.ctx, "auth.usecase.Logout")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.AccessToken).
					Return(nil, assert.AnError)

				return &Logout{
					tel:       tel,
					validator: validatorMock,
					secHash:   secHashMock,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreDeleteTokenByAccess",
			args: args{
				ctx: context.Background(),
				in:  domain.LogoutInput{AccessToken: "token"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Logout {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockLogoutStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Logout")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.AccessToken).
					Return([]byte("hash"), nil)

				storeMock.EXPECT().
					DeleteTokenByAccess(ctx, "hash").
					Return(assert.AnError)

				return &Logout{
					tel:       tel,
					validator: validatorMock,
					secHash:   secHashMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  domain.LogoutInput{AccessToken: "token"},
			},
			want:    &domain.LogoutOutput{Message: "You have successfully logged out!"},
			wantErr: nil,
			mockFn: func(a args) *Logout {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				secHashMock := mockHash.NewMockHash(t)
				storeMock := mockz.NewMockLogoutStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "auth.usecase.Logout")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				secHashMock.EXPECT().
					Hash(a.in.AccessToken).
					Return([]byte("hash"), nil)

				storeMock.EXPECT().
					DeleteTokenByAccess(ctx, "hash").
					Return(nil)

				return &Logout{
					tel:       tel,
					validator: validatorMock,
					secHash:   secHashMock,
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
