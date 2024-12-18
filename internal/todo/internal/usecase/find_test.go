package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	vm "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewFind(t *testing.T) {
	type args struct {
		dep Dependency
		s   FindStore
	}
	tests := []struct {
		name string
		args args
		want *Find
	}{
		{
			name: "Success",
			args: args{},
			want: &Find{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewFind(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind_Execute(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.FindInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Todo
		wantErr error
		mockFn  func(a args) *Find
	}{
		{
			name: "ErrorValidation",
			args: args{
				ctx: context.Background(),
				in:  domain.FindInput{ID: 11},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Find {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				_, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Find")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Find{
					telemetry: mtel,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name: "ErrorStore",
			args: args{
				ctx: context.Background(),
				in:  domain.FindInput{ID: 11},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Find {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockFindStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Find")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				store.EXPECT().
					Find(ctx, a.in.ID).
					Return(nil, assert.AnError)

				return &Find{
					telemetry: mtel,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name: "StoreNotFound",
			args: args{
				ctx: context.Background(),
				in:  domain.FindInput{ID: 11},
			},
			want:    nil,
			wantErr: goerror.NewBusiness("todo not found", goerror.CodeNotFound),
			mockFn: func(a args) *Find {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockFindStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Find")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				store.EXPECT().
					Find(ctx, a.in.ID).
					Return(nil, nil)

				return &Find{
					telemetry: mtel,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in:  domain.FindInput{ID: 10},
			},
			want: &domain.Todo{
				ID:          10,
				UserID:      11,
				Title:       "test 1",
				Description: "test 2",
				Status:      enum.New(domain.TodoStatusDrop),
			},
			wantErr: nil,
			mockFn: func(a args) *Find {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockFindStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Find")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				todo := &domain.Todo{
					ID:          10,
					UserID:      11,
					Title:       "test 1",
					Description: "test 2",
					Status:      enum.New(domain.TodoStatusDrop),
				}
				store.EXPECT().
					Find(ctx, a.in.ID).
					Return(todo, nil)

				return &Find{
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
