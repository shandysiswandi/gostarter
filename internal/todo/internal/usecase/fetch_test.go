package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/enum"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/pagination"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestNewFetch(t *testing.T) {
	type args struct {
		dep Dependency
		s   FetchStore
	}
	tests := []struct {
		name string
		args args
		want *Fetch
	}{
		{
			name: "Success",
			args: args{},
			want: &Fetch{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewFetch(tt.args.dep, tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFetch_Execute(t *testing.T) {
	type args struct {
		ctx context.Context
		in  domain.FetchInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.FetchOutput
		wantErr error
		mockFn  func(a args) *Fetch
	}{
		{
			name: "ErrorStore",
			args: args{
				ctx: context.Background(),
				in: domain.FetchInput{
					Cursor: "NTY",
					Limit:  "1",
					Status: "",
				},
			},
			want:    nil,
			wantErr: goerror.NewServer("failed to fetch todo", assert.AnError),
			mockFn: func(a args) *Fetch {
				mtel := telemetry.NewTelemetry()
				store := mockz.NewMockFetchStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Fetch")
				defer span.End()

				cursor, limit := pagination.ParseCursorBased(a.in.Cursor, a.in.Limit)

				filter := map[string]any{
					"cursor": cursor,
					"limit":  limit,
				}
				store.EXPECT().
					Fetch(ctx, filter).
					Return(nil, assert.AnError)

				return &Fetch{
					telemetry: mtel,
					store:     store,
				}
			},
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				in: domain.FetchInput{
					Cursor: "",
					Limit:  "1",
					Status: "done",
				},
			},
			want: &domain.FetchOutput{
				Todos: []domain.Todo{{
					ID:          1,
					UserID:      2,
					Title:       "test 1",
					Description: "test 1",
					Status:      enum.New(domain.TodoStatusDone),
				}},
				NextCursor: "Mg",
				HasMore:    true,
			},
			wantErr: nil,
			mockFn: func(a args) *Fetch {
				mtel := telemetry.NewTelemetry()
				store := mockz.NewMockFetchStore(t)

				ctx, span := mtel.Tracer().Start(a.ctx, "todo.usecase.Fetch")
				defer span.End()

				cursor, limit := pagination.ParseCursorBased(a.in.Cursor, a.in.Limit)

				filter := map[string]any{
					"cursor": cursor,
					"limit":  limit,
					"status": enum.New(enum.Parse[domain.TodoStatus](a.in.Status)),
				}
				todos := []domain.Todo{
					{
						ID:          1,
						UserID:      2,
						Title:       "test 1",
						Description: "test 1",
						Status:      enum.New(domain.TodoStatusDone),
					},
					{
						ID:          2,
						UserID:      2,
						Title:       "test 2",
						Description: "test 2",
						Status:      enum.New(domain.TodoStatusDone),
					},
				}
				store.EXPECT().
					Fetch(ctx, filter).
					Return(todos, nil)

				return &Fetch{
					telemetry: mtel,
					store:     store,
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
