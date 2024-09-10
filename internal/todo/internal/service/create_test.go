package service

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/logger"
	lm "github.com/shandysiswandi/gostarter/pkg/logger/mocker"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	um "github.com/shandysiswandi/gostarter/pkg/uid/mocker"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	vm "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewCreate(t *testing.T) {
	type args struct {
		l     logger.Logger
		s     CreateStore
		v     validation.Validator
		idgen uid.NumberID
	}
	tests := []struct {
		name string
		args args
		want *Create
	}{
		{name: "Success", args: args{}, want: &Create{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCreate(tt.args.l, tt.args.s, tt.args.v, tt.args.idgen)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreate_Execute(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.CreateInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.CreateOutput
		wantErr error
		mockFn  func(a args) *Create
	}{
		{
			name:    "ErrorValidation",
			args:    args{ctx: context.TODO(), in: domain.CreateInput{}},
			want:    nil,
			wantErr: goerror.NewInvalidInput("validation input fail", assert.AnError),
			mockFn: func(a args) *Create {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)
				log.EXPECT().Warn(a.ctx, "validation failed").Return()

				return &Create{
					log:       log,
					store:     nil,
					uidnumber: nil,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorNotAffected",
			args:    args{ctx: context.TODO(), in: domain.CreateInput{}},
			want:    nil,
			wantErr: goerror.NewBusiness("failed to create todo", goerror.CodeUnknown),
			mockFn: func(a args) *Create {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				idgen := um.NewMockNumberID(t)
				store := mockz.NewMockCreateStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				idgen.EXPECT().Generate().Return(101)

				input := domain.Todo{
					ID:          101,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      domain.TodoStatusInitiate,
				}
				store.EXPECT().Create(a.ctx, input).Return(domain.ErrTodoNotCreated)
				log.EXPECT().Warn(a.ctx, "todo created but db not affected").Return()

				return &Create{
					log:       log,
					store:     store,
					uidnumber: idgen,
					validator: validator,
				}
			},
		},
		{
			name:    "ErrorStore",
			args:    args{ctx: context.TODO(), in: domain.CreateInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to create todo", assert.AnError),
			mockFn: func(a args) *Create {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				idgen := um.NewMockNumberID(t)
				store := mockz.NewMockCreateStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				idgen.EXPECT().Generate().Return(101)

				input := domain.Todo{
					ID:          101,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      domain.TodoStatusInitiate,
				}
				store.EXPECT().Create(a.ctx, input).Return(assert.AnError)
				log.EXPECT().Error(a.ctx, "todo fail to create", assert.AnError).Return()

				return &Create{
					log:       log,
					store:     store,
					uidnumber: idgen,
					validator: validator,
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), in: domain.CreateInput{}},
			want:    &domain.CreateOutput{ID: 101},
			wantErr: nil,
			mockFn: func(a args) *Create {
				log := lm.NewMockLogger(t)
				validator := vm.NewMockValidator(t)
				idgen := um.NewMockNumberID(t)
				store := mockz.NewMockCreateStore(t)

				validator.EXPECT().Validate(a.in).Return(nil)

				idgen.EXPECT().Generate().Return(101)

				input := domain.Todo{
					ID:          101,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      domain.TodoStatusInitiate,
				}
				store.EXPECT().Create(a.ctx, input).Return(nil)

				return &Create{
					log:       log,
					store:     store,
					uidnumber: idgen,
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
