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

func TestNewDelete(t *testing.T) {
	type args struct {
		dep Dependency
		s   DeleteStore
	}
	tests := []struct {
		name string
		args args
		want *Delete
	}{
		{
			name: "Success",
			args: args{},
			want: &Delete{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewDelete(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDelete_Call(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.DeleteInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.DeleteOutput
		wantErr error
		mockFn  func(a args) *Delete
	}{
		{
			name:    "ErrorValidation",
			args:    args{ctx: context.TODO(), in: domain.DeleteInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Delete {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)

				return &Delete{
					telemetry: mtel,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorStore",
			args:    args{ctx: context.TODO(), in: domain.DeleteInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to delete todo", assert.AnError),
			mockFn: func(a args) *Delete {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockDeleteStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Delete(a.ctx, a.in.ID).Return(assert.AnError)

				return &Delete{
					telemetry: mtel,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), in: domain.DeleteInput{ID: 111}},
			want:    &domain.DeleteOutput{ID: 111},
			wantErr: nil,
			mockFn: func(a args) *Delete {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockDeleteStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Delete(a.ctx, a.in.ID).Return(nil)

				return &Delete{
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
