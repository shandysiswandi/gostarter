package service

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	vm "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewUpdate(t *testing.T) {
	type args struct {
		t *telemetry.Telemetry
		s UpdateStore
		v validation.Validator
	}
	tests := []struct {
		name string
		args args
		want *Update
	}{
		{name: "Success", args: args{}, want: &Update{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewUpdate(tt.args.t, tt.args.s, tt.args.v)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUpdate_Execute(t *testing.T) {
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
			name:    "ErrorValidation",
			args:    args{ctx: context.TODO(), in: domain.UpdateInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Update {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)

				return &Update{
					telemetry: mtel,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorStore",
			args:    args{ctx: context.TODO(), in: domain.UpdateInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to update todo", assert.AnError),
			mockFn: func(a args) *Update {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				sts := domain.ParseTodoStatus(a.in.Status)
				store.EXPECT().Update(a.ctx, domain.Todo{
					ID:          a.in.ID,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      sts,
				}).Return(assert.AnError)

				return &Update{
					telemetry: mtel,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name: "Success",
			args: args{ctx: context.TODO(), in: domain.UpdateInput{
				ID:          120,
				Title:       "test 1",
				Description: "test 2",
				Status:      "DONE",
			}},
			want: &domain.Todo{
				ID:          120,
				Title:       "test 1",
				Description: "test 2",
				Status:      domain.TodoStatusDone,
			},
			wantErr: nil,
			mockFn: func(a args) *Update {
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockUpdateStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				sts := domain.ParseTodoStatus(a.in.Status)
				store.EXPECT().Update(a.ctx, domain.Todo{
					ID:          a.in.ID,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      sts,
				}).Return(nil)

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
			got, err := s.Execute(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
