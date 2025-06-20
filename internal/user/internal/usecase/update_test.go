package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/goerror"
	mockValidation "github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/user/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/user/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdate(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		s    UpdateStore
		want *Update
	}{
		{
			name: "Success",
			dep:  Dependency{},
			s:    nil,
			want: &Update{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUpdate(tt.dep, tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdate_Call(t *testing.T) {
	claim := lib.NewJWTClaim(11, "email", time.Time{}, nil)
	ctxJWT := lib.SetJWTClaim(context.Background(), claim)

	type args struct {
		ctx context.Context
		in  domain.UpdateInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr error
		mockFn  func(a args) *Update
	}{
		{
			name: "ErrorValidationInput",
			args: args{
				ctx: ctxJWT,
				in:  domain.UpdateInput{Name: "name"},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *Update {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)

				_, span := tel.Tracer().Start(a.ctx, "user.usecase.Update")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Update{
					tel:       tel,
					validator: validatorMock,
					store:     nil,
				}
			},
		},
		{
			name: "ErrorStoreUpdate",
			args: args{
				ctx: ctxJWT,
				in:  domain.UpdateInput{Name: "full name"},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Update {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockUpdateStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.Update")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				in := map[string]any{
					"id":   uint64(11),
					"name": a.in.Name,
				}
				storeMock.EXPECT().
					Update(ctx, in).
					Return(assert.AnError)

				return &Update{
					tel:       tel,
					validator: validatorMock,
					store:     storeMock,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: ctxJWT,
				in:  domain.UpdateInput{Name: "full name"},
			},
			want: &domain.User{
				ID:       11,
				Name:     "full name",
				Email:    "email",
				Password: "***",
			},
			wantErr: nil,
			mockFn: func(a args) *Update {
				tel := telemetry.NewTelemetry()
				validatorMock := mockValidation.NewMockValidator(t)
				storeMock := mockz.NewMockUpdateStore(t)

				ctx, span := tel.Tracer().Start(a.ctx, "user.usecase.Update")
				defer span.End()

				validatorMock.EXPECT().
					Validate(a.in).
					Return(nil)

				in := map[string]any{
					"id":   uint64(11),
					"name": a.in.Name,
				}
				storeMock.EXPECT().
					Update(ctx, in).
					Return(nil)

				return &Update{
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
