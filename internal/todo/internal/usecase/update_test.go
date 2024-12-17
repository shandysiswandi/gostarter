package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/jwt"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	vm "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdate(t *testing.T) {
	type args struct {
		dep Dependency
		s   UpdateStore
	}
	tests := []struct {
		name string
		args args
		want *Update
	}{
		{
			name: "Success",
			args: args{},
			want: &Update{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUpdate(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdate_Execute(t *testing.T) {
	claim := jwt.NewClaim(11, "email", time.Time{}, nil)
	ctx := jwt.SetClaim(context.Background(), claim)

	type args struct {
		ctx context.Context
		in  domain.UpdateInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Todo
		wantErr error
		mockFn  func(a args) *Update
	}{
		{
			name: "ErrorValidation",
			args: args{
				ctx: ctx,
				in: domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "done",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Update {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				_, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Update")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Update{
					telemetry: mtel,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name: "ErrorStore",
			args: args{
				ctx: ctx,
				in: domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "done",
				},
			},
			want:    nil,
			wantErr: goerror.NewServer("failed to update todo", assert.AnError),
			mockFn: func(a args) *Update {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Update")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				sts := enum.New(enum.Parse[domain.TodoStatus](a.in.Status))
				data := domain.Todo{
					ID:          a.in.ID,
					UserID:      11,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      sts,
				}
				store.EXPECT().
					Update(ctx, data).
					Return(assert.AnError)

				return &Update{
					telemetry: mtel,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: ctx,
				in: domain.UpdateInput{
					ID:          10,
					Title:       "title",
					Description: "description",
					Status:      "DONE",
				},
			},
			want: &domain.Todo{
				ID:          10,
				UserID:      11,
				Title:       "title",
				Description: "description",
				Status:      enum.New(domain.TodoStatusDone),
			},
			wantErr: nil,
			mockFn: func(a args) *Update {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Update")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				data := domain.Todo{
					ID:          a.in.ID,
					UserID:      11,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      enum.New(enum.Parse[domain.TodoStatus](a.in.Status)),
				}
				store.EXPECT().
					Update(ctx, data).
					Return(nil)

				return &Update{
					telemetry: mtel,
					store:     store,
					validator: validator,
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
