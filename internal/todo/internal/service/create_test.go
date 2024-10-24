package service

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/shandysiswandi/gostarter/pkg/uid"
	um "github.com/shandysiswandi/gostarter/pkg/uid/mocker"
	"github.com/shandysiswandi/gostarter/pkg/validation"
	vm "github.com/shandysiswandi/gostarter/pkg/validation/mocker"
	"github.com/stretchr/testify/assert"
)

func TestNewCreate(t *testing.T) {
	type args struct {
		t     *telemetry.Telemetry
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
			got := NewCreate(tt.args.t, tt.args.s, tt.args.v, tt.args.idgen)
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
				mtel := telemetry.NewTelemetry()
				validator := vm.NewMockValidator(t)

				validator.EXPECT().Validate(a.in).Return(assert.AnError)

				return &Create{
					telemetry: mtel,
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
				mtel := telemetry.NewTelemetry()
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

				return &Create{
					telemetry: mtel,
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
				mtel := telemetry.NewTelemetry()
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

				return &Create{
					telemetry: mtel,
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
				mtel := telemetry.NewTelemetry()
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
					telemetry: mtel,
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
