package usecase

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/mockz"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
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
		want    []domain.Todo
		wantErr error
		mockFn  func(a args) *Fetch
	}{
		{
			name:    "ErrorStore",
			args:    args{ctx: context.TODO(), in: domain.FetchInput{}},
			want:    nil,
			wantErr: goerror.NewServer("failed to fetch todo", assert.AnError),
			mockFn: func(a args) *Fetch {
				mtel := telemetry.NewTelemetry()
				store := mockz.NewMockFetchStore(t)

				filter := map[string]string{
					"id":          a.in.ID,
					"title":       a.in.Title,
					"description": a.in.Description,
					"status":      a.in.Status,
				}
				store.EXPECT().Fetch(a.ctx, filter).Return(nil, assert.AnError)

				return &Fetch{
					telemetry: mtel,
					store:     store,
				}
			},
		},
		{
			name: "Success",
			args: args{ctx: context.TODO(), in: domain.FetchInput{}},
			want: []domain.Todo{{
				ID:          91,
				Title:       "test 1",
				Description: "test 2",
				Status:      domain.TodoStatusInitiate,
			}},
			wantErr: nil,
			mockFn: func(a args) *Fetch {
				mtel := telemetry.NewTelemetry()
				store := mockz.NewMockFetchStore(t)

				filter := map[string]string{
					"id":          a.in.ID,
					"title":       a.in.Title,
					"description": a.in.Description,
					"status":      a.in.Status,
				}
				store.EXPECT().Fetch(a.ctx, filter).Return([]domain.Todo{{
					ID:          91,
					Title:       "test 1",
					Description: "test 2",
					Status:      domain.TodoStatusInitiate,
				}}, nil)

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
