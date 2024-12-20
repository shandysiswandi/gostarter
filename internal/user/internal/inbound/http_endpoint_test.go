package inbound

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
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
				c := framework.NewTestContext(http.MethodGet, "/me/profile", nil)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				profileMock := mockz.NewMockProfile(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Profile")
				defer span.End()

				in := domain.ProfileInput{}
				profileMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:       tel,
					profileUC: profileMock,
					updateUC:  nil,
					logoutUC:  nil,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodGet, "/me/profile", nil)
				return c.Build()
			},
			want: User{
				ID:    1,
				Name:  "full name",
				Email: "email",
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				profileMock := mockz.NewMockProfile(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Profile")
				defer span.End()

				in := domain.ProfileInput{}
				out := &domain.User{
					ID:       1,
					Name:     "full name",
					Email:    "email",
					Password: "***",
				}
				profileMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:       tel,
					profileUC: profileMock,
					updateUC:  nil,
					logoutUC:  nil,
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

func Test_httpEndpoint_Update(t *testing.T) {
	tests := []struct {
		name    string
		c       func() framework.Context
		want    any
		wantErr error
		mockFn  func(ctx context.Context) *httpEndpoint
	}{
		{
			name: "ErrorDecodeBody",
			c: func() framework.Context {
				body := bytes.NewBufferString("fake request")
				c := framework.NewTestContext(http.MethodPatch, "/me/profile", body)
				return c.Build()
			},
			want:    nil,
			wantErr: goerror.NewInvalidFormat("invalid request body"),
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Update")
				defer span.End()

				return &httpEndpoint{
					tel:       tel,
					profileUC: nil,
					updateUC:  updateMock,
					logoutUC:  nil,
				}
			},
		},
		{
			name: "ErrorCallUC",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"name":"full name"}`)
				c := framework.NewTestContext(http.MethodPatch, "/me/profile", body)
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{Name: "full name"}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:       tel,
					profileUC: nil,
					updateUC:  updateMock,
					logoutUC:  nil,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				body := bytes.NewBufferString(`{"name":"full name"}`)
				c := framework.NewTestContext(http.MethodPatch, "/me/profile", body)
				return c.Build()
			},
			want: User{
				ID:    21,
				Name:  "full name",
				Email: "full@name.com",
			},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				updateMock := mockz.NewMockUpdate(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Update")
				defer span.End()

				in := domain.UpdateInput{Name: "full name"}
				out := &domain.User{
					ID:       21,
					Name:     "full name",
					Email:    "full@name.com",
					Password: "***",
				}
				updateMock.EXPECT().
					Call(ctx, in).
					Return(out, nil)

				return &httpEndpoint{
					tel:       tel,
					profileUC: nil,
					updateUC:  updateMock,
					logoutUC:  nil,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := tt.c()
			e := tt.mockFn(c.Context())
			got, err := e.Update(c)
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
				c := framework.NewTestContext(http.MethodPost, "/me/logout", nil)
				c.SetHeader("Authorization", "Bearer ay")
				return c.Build()
			},
			want:    nil,
			wantErr: assert.AnError,
			mockFn: func(ctx context.Context) *httpEndpoint {
				logoutMock := mockz.NewMockLogout(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Logout")
				defer span.End()

				in := domain.LogoutInput{AccessToken: "Bearer ay"}
				logoutMock.EXPECT().
					Call(ctx, in).
					Return(nil, assert.AnError)

				return &httpEndpoint{
					tel:       tel,
					profileUC: nil,
					updateUC:  nil,
					logoutUC:  logoutMock,
				}
			},
		},
		{
			name: "Success",
			c: func() framework.Context {
				c := framework.NewTestContext(http.MethodPost, "/me/logout", nil)
				c.SetHeader("Authorization", "Bearer ay")
				return c.Build()
			},
			want:    LogoutResponse{Message: "success"},
			wantErr: nil,
			mockFn: func(ctx context.Context) *httpEndpoint {
				logoutMock := mockz.NewMockLogout(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(ctx, "user.inbound.httpEndpoint.Logout")
				defer span.End()

				in := domain.LogoutInput{AccessToken: "Bearer ay"}
				logoutMock.EXPECT().
					Call(ctx, in).
					Return(&domain.LogoutOutput{Message: "success"}, nil)

				return &httpEndpoint{
					tel:       tel,
					profileUC: nil,
					updateUC:  nil,
					logoutUC:  logoutMock,
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
