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

func TestNewUpdateStatus(t *testing.T) {
	type args struct {
		dep Dependency
		s   UpdateStatusStore
	}
	tests := []struct {
		name string
		args args
		want *UpdateStatus
	}{
		{
			name: "Success",
			args: args{},
			want: &UpdateStatus{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUpdateStatus(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdateStatus_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.UpdateStatusInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.UpdateStatusOutput
		wantErr error
		mockFn  func(a args) *UpdateStatus
	}{
		{
			name: "ErrorValidation",
			args: args{
				ctx: context.Background(),
				in: domain.UpdateStatusInput{
					ID:     10,
					Status: "done",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *UpdateStatus {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				_, span := mtel.Tracer().Start(a.ctx, "todo.usecase.UpdateStatus")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &UpdateStatus{
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
				in: domain.UpdateStatusInput{
					ID:     10,
					Status: "done",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *UpdateStatus {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStatusStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.UpdateStatus")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				sts := enum.New(enum.Parse[domain.TodoStatus](a.in.Status))
				store.EXPECT().
					UpdateStatus(ctx, a.in.ID, sts).
					Return(assert.AnError)

				return &UpdateStatus{
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
				in: domain.UpdateStatusInput{
					ID:     10,
					Status: "done",
				},
			},
			want: &domain.UpdateStatusOutput{
				ID:     10,
				Status: enum.New(domain.TodoStatusUnknown),
			},
			wantErr: nil,
			mockFn: func(a args) *UpdateStatus {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStatusStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.UpdateStatus")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				sts := enum.New(enum.Parse[domain.TodoStatus](a.in.Status))
				store.EXPECT().
					UpdateStatus(ctx, a.in.ID, sts).
					Return(nil)

				return &UpdateStatus{
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
