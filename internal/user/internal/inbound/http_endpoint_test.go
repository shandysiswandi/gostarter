package inbound

import (
	"context"
	"net/http"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func Test_httpEndpoint_Profile(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/users/profile", nil)
				ctx := jwt.SetClaim(context.Background(), &jwt.Claim{Email: "email"})
				c.SetContext(ctx)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				profileMock := mockz.NewMockProfile(t)

				in := domain.ProfileInput{Email: "email"}
				profileMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					profileUC: profileMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/users/profile", nil)
				ctx := jwt.SetClaim(context.Background(), &jwt.Claim{Email: "email"})
				c.SetContext(ctx)
				return c.Build()
			},
			want:    ProfileResponse{ID: 1, Email: "email"},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				profileMock := mockz.NewMockProfile(t)

				in := domain.ProfileInput{Email: "email"}
				profileMock.EXPECT().
					Call(ctx, in).
					Return(&domain.User{ID: 1, Email: "email"}, nil)

				return &httpEndpoint{
					profileUC: profileMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Profile(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_httpEndpoint_Logout(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPost, "/users/logout", nil)
				c.SetHeader("Authorization", "Bearer ay")
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				logoutMock := mockz.NewMockLogout(t)

				in := domain.LogoutInput{AccessToken: "Bearer ay"}
				logoutMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					logoutUC: logoutMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPost, "/users/logout", nil)
				c.SetHeader("Authorization", "Bearer ay")
				return c.Build()
			},
			want:    LogoutResponse{Message: "success"},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				logoutMock := mockz.NewMockLogout(t)

				in := domain.LogoutInput{AccessToken: "Bearer ay"}
				logoutMock.EXPECT().
					Call(ctx, in).
					Return(&domain.LogoutOutput{Message: "success"}, nil)

				return &httpEndpoint{
					logoutUC: logoutMock,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Logout(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
