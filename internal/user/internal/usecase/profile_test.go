package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
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
	claim := lib.NewJWTClaim(11, "email", time.Time{}, nil)
	ctxJWT := lib.SetJWTClaim(context.Background(), claim)

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
			name: "ErrorStoreFindUserByEmail",
			args: args{
				ctx: ctxJWT,
				in:  domain.ProfileInput{},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				storeMock := mockz.NewMockProfileStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.Profile")
				defer span.End()

				storeMock.EXPECT().
					FindUserByEmail(ctx, "email").
					Return(nil, assert.AnError)

				return &Profile{
					tel:   tel,
					store: storeMock,
				}
			},
		},
		{
			name: "ErrorStoreFindUserByEmailNotFound",
			args: args{
				ctx: ctxJWT,
				in:  domain.ProfileInput{},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("user not found", goerror.CodeNotFound),
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				storeMock := mockz.NewMockProfileStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.Profile")
				defer span.End()

				storeMock.EXPECT().
					FindUserByEmail(ctx, "email").
					Return(nil, nil)

				return &Profile{
					tel:   tel,
					store: storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: ctxJWT,
				in:  domain.ProfileInput{},
			},
			want: &domain.User{
				ID:       1,
				Name:     "full name",
				Email:    "email",
				Password: "***",
			},
			wantErr: nil,
			mockFn: func(a args) *Profile {
				tel := telemetry.NewTelemetry()
				storeMock := mockz.NewMockProfileStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.Profile")
				defer span.End()

				out := &domain.User{
					ID:       1,
					Name:     "full name",
					Email:    "email",
					Password: "password",
				}
				storeMock.EXPECT().
					FindUserByEmail(ctx, "email").
					Return(out, nil)

				return &Profile{
					tel:   tel,
					store: storeMock,
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
