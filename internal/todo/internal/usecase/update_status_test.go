package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
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

func TestUpdateStatus_Execute(t *testing.T) {
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
			name:    "ErrorValidation",
			args:    args{ctx: context.Background(), in: domain.UpdateStatusInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *UpdateStatus {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)

				return &UpdateStatus{
					telemetry: mtel,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorStore",
			args:    args{ctx: context.Background(), in: domain.UpdateStatusInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to update status todo", assert.AnError),
			mockFn: func(a args) *UpdateStatus {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStatusStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				sts := domain.ParseTodoStatus(a.in.Status)
				store.EXPECT().UpdateStatus(a.ctx, a.in.ID, sts).Return(assert.AnError)

				return &UpdateStatus{
					telemetry: mtel,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.Background(), in: domain.UpdateStatusInput{ID: 1}},
			want:    &domain.UpdateStatusOutput{ID: 1, Status: domain.TodoStatusUnknown},
			wantErr: nil,
			mockFn: func(a args) *UpdateStatus {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStatusStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				sts := domain.ParseTodoStatus(a.in.Status)
				store.EXPECT().UpdateStatus(a.ctx, a.in.ID, sts).Return(nil)

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
