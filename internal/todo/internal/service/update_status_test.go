package service

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	lm "github.com/shandysiswandi/gostarter/pkg/logger/mocker"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	vm "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdateStatus(t *testing.T) {
	type args struct {
		l logger.Logger
		s UpdateStatusStore
		v validation.Validator
	}
	tests := []struct {
		name string
		args args
		want *UpdateStatus
	}{
		{name: "Success", args: args{}, want: &UpdateStatus{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUpdateStatus(tt.args.l, tt.args.s, tt.args.v)
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
			args:    args{ctx: context.TODO(), in: domain.UpdateStatusInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *UpdateStatus {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)
				log.EXPECT().Warn(a.ctx, "validation failed").Return()

				return &UpdateStatus{
					log:       log,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorStore",
			args:    args{ctx: context.TODO(), in: domain.UpdateStatusInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to update status todo", assert.AnError),
			mockFn: func(a args) *UpdateStatus {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStatusStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				sts := domain.ParseTodoStatus(a.in.Status)
				store.EXPECT().UpdateStatus(a.ctx, a.in.ID, sts).Return(assert.AnError)
				log.EXPECT().Error(a.ctx, "todo fail to update status", assert.AnError).Return()

				return &UpdateStatus{
					log:       log,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), in: domain.UpdateStatusInput{ID: 1}},
			want:    &domain.UpdateStatusOutput{ID: 1, Status: domain.TodoStatusUnknown},
			wantErr: nil,
			mockFn: func(a args) *UpdateStatus {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStatusStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				sts := domain.ParseTodoStatus(a.in.Status)
				store.EXPECT().UpdateStatus(a.ctx, a.in.ID, sts).Return(nil)

				return &UpdateStatus{
					log:       log,
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
			got, err := s.Execute(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
