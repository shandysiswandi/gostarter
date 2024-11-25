package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	mockValidation "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewProfile(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		s    ProfileStore
		want *Profile
	}{
		{
			name: "Success",
			dep:  Dependency{},
			s:    nil,
			want: &Profile{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewProfile(tt.dep, tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProfile_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.ProfileInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr error
		mockFn  func(a args) *Profile
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: context.Background(),
				in:  domain.ProfileInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "usecase.Profile")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Profile{
					tel:       tel,
					validator: validatorMock,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmail",
			args: args{
				ctx: context.Background(),
				in:  domain.ProfileInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockProfileStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "usecase.Profile")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, assert.AnError)

				return &Profile{
					tel:       tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmailNotFound",
			args: args{
				ctx: context.Background(),
				in:  domain.ProfileInput{Email: "email"},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("user not found", goerror.CodeNotFound),
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockProfileStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "usecase.Profile")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(nil, nil)

				return &Profile{
					tel:       tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  domain.ProfileInput{Email: "email"},
			},
			want:    &domain.User{ID: 1, Email: "email", Password: "***"},
			wantErr: nil,
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockProfileStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "usecase.Profile")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				storeMock.EXPECT().
					FindUserByEmail(ctx, a.in.Email).
					Return(&domain.User{ID: 1, Email: "email", Password: "password"}, nil)

				return &Profile{
					tel:       tel,
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
