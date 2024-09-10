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

func TestNewFind(t *testing.T) {
	type args struct {
		l logger.Logger
		s FindStore
		v validation.Validator
	}
	tests := []struct {
		name string
		args args
		want *Find
	}{
		{name: "Success", args: args{}, want: &Find{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewFind(tt.args.l, tt.args.s, tt.args.v)
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
			name:    "ErrorValidation",
			args:    args{ctx: context.TODO(), in: domain.FindInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Find {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)
				log.EXPECT().Warn(a.ctx, "validation failed").Return()

				return &Find{
					log:       log,
					store:     nil,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorStore",
			args:    args{ctx: context.TODO(), in: domain.FindInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to find todo", assert.AnError),
			mockFn: func(a args) *Find {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockFindStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Find(a.ctx, a.in.ID).Return(nil, assert.AnError)
				log.EXPECT().Error(a.ctx, "todo fail to find", assert.AnError).Return()

				return &Find{
					log:       log,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name:    "StoreNotFound",
			args:    args{ctx: context.TODO(), in: domain.FindInput{}},
			want:    nil,
			wantErr: goerror.NewBusiness("todo not found", goerror.CodeNotFound),
			mockFn: func(a args) *Find {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockFindStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Find(a.ctx, a.in.ID).Return(nil, nil)
				log.EXPECT().Warn(a.ctx, "todo is not found").Return()

				return &Find{
					log:       log,
					store:     store,
					validator: validator,
				}
			},
		},
		{
			name: "Success",
			args: args{ctx: context.TODO(), in: domain.FindInput{}},
			want: &domain.Todo{
				ID:          10,
				Title:       "test 1",
				Description: "test 2",
				Status:      domain.TodoStatusDrop,
			},
			wantErr: nil,
			mockFn: func(a args) *Find {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				store := mockz.NewMockFindStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				store.EXPECT().Find(a.ctx, a.in.ID).Return(&domain.Todo{
					ID:          10,
					Title:       "test 1",
					Description: "test 2",
					Status:      domain.TodoStatusDrop,
				}, nil)

				return &Find{
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
