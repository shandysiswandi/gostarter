package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/shandysiswandi/goreng/enum"
	"github.com/shandysiswandi/goreng/goerror"
	um "github.com/shandysiswandi/goreng/mocker"
	vm "github.com/shandysiswandi/goreng/mocker"
	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/lib"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/stretchr/testify/assert"
)

func TestNewCreate(t *testing.T) {
	type args struct {
		dep Dependency
		s   CreateStore
	}
	tests := []struct {
		name string
		args args
		want *Create
	}{
		{
			name: "Success",
			args: args{},
			want: &Create{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewCreate(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreate_Execute(t *testing.T) {
	claim := lib.NewJWTClaim(11, "email", time.Time{}, nil)
	ctx := lib.SetJWTClaim(context.Background(), claim)

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
			name: "ErrorValidation",
			args: args{
				ctx: ctx,
				in: domain.CreateInput{
					Title:       "title",
					Description: "description",
				},
			},
			want:    nil,
			wantErr: goerror.NewInvalidInput("Invalid request payload", assert.AnError),
			mockFn: func(a args) *Create {
				validator := vm.NewMockValidator(t)
				tel := telemetry.NewTelemetry()

				_, span := tel.Tracer().Start(a.ctx, "todo.usecase.Create")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(assert.AnError)

				return &Create{
					telemetry: tel,
					store:     nil,
					uidnumber: nil,
					validator: validator,
				}
			},
		},
		{
			name: "ErrorNotAffected",
			args: args{
				ctx: ctx,
				in: domain.CreateInput{
					Title:       "title",
					Description: "description",
				}},
			want:    nil,
			wantErr: goerror.NewBusiness("failed to create todo", goerror.CodeUnknown),
			mockFn: func(a args) *Create {
				validator := vm.NewMockValidator(t)
				idgen := um.NewMockNumberID(t)
				store := mockz.NewMockCreateStore(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.usecase.Create")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				idgen.EXPECT().
					Generate().
					Return(101)

				input := domain.Todo{
					ID:          101,
					UserID:      11,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      enum.New(domain.TodoStatusInitiate),
				}
				store.EXPECT().
					Create(ctx, input).
					Return(domain.ErrTodoNotCreated)

				return &Create{
					telemetry: tel,
					store:     store,
					uidnumber: idgen,
					validator: validator,
				}
			},
		},
		{
			name: "ErrorStore",
			args: args{
				ctx: ctx,
				in: domain.CreateInput{
					Title:       "title",
					Description: "description",
				},
			},
			want:    nil,
			wantErr: goerror.NewServerInternal(assert.AnError),
			mockFn: func(a args) *Create {
				validator := vm.NewMockValidator(t)
				idgen := um.NewMockNumberID(t)
				store := mockz.NewMockCreateStore(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.usecase.Create")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				idgen.EXPECT().
					Generate().
					Return(101)

				input := domain.Todo{
					ID:          101,
					UserID:      11,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      enum.New(domain.TodoStatusInitiate),
				}
				store.EXPECT().
					Create(ctx, input).
					Return(assert.AnError)

				return &Create{
					telemetry: tel,
					store:     store,
					uidnumber: idgen,
					validator: validator,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: ctx,
				in: domain.CreateInput{
					Title:       "title",
					Description: "description",
				},
			},
			want:    &domain.CreateOutput{ID: 101},
			wantErr: nil,
			mockFn: func(a args) *Create {
				validator := vm.NewMockValidator(t)
				idgen := um.NewMockNumberID(t)
				store := mockz.NewMockCreateStore(t)
				tel := telemetry.NewTelemetry()

				ctx, span := tel.Tracer().Start(a.ctx, "todo.usecase.Create")
				defer span.End()

				validator.EXPECT().
					Validate(a.in).
					Return(nil)

				idgen.EXPECT().
					Generate().
					Return(101)

				input := domain.Todo{
					ID:          101,
					UserID:      11,
					Title:       a.in.Title,
					Description: a.in.Description,
					Status:      enum.New(domain.TodoStatusInitiate),
				}
				store.EXPECT().
					Create(ctx, input).
					Return(nil)

				return &Create{
					telemetry: tel,
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
			got, err := s.Call(tt.args.ctx, tt.args.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
