package goroutine

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	type args struct {
		maxGoroutine int
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		{
			name: "Success",
			args: args{
				maxGoroutine: 2,
			},
			want: &Manager{
				wg:   &sync.WaitGroup{},
				sema: make(chan struct{}, 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewManager(tt.args.maxGoroutine)
			assert.Equal(t, tt.want.errs, got.errs)
			assert.Equal(t, tt.want.wg, got.wg)
		})
	}
}

func TestManager_Go(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		ctx context.Context
		f   func(c context.Context) error
	}
	tests := []struct {
		name   string
		args   args
		mockFn func(a args) *Manager
	}{
		{
			name: "Default",
			args: args{
				ctx: context.Background(),
				f: func(c context.Context) error {
					return nil
				},
			},
			mockFn: func(a args) *Manager {
				return NewManager(0)
			},
		},
		{
			name: "FuncReturnError",
			args: args{
				ctx: context.Background(),
				f: func(c context.Context) error {
					return assert.AnError
				},
			},
			mockFn: func(a args) *Manager {
				return NewManager(1)
			},
		},
		{
			name: "FuncHavePanic",
			args: args{
				ctx: context.Background(),
				f: func(c context.Context) error {
					panic(1)
				},
			},
			mockFn: func(a args) *Manager {
				return NewManager(1)
			},
		},
		{
			name: "ContextDone",
			args: args{
				ctx: ctx,
				f: func(c context.Context) error {
					return nil
				},
			},
			mockFn: func(a args) *Manager {
				return NewManager(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := tt.mockFn(tt.args)

			g.Go(tt.args.ctx, tt.args.f)

			g.Wait()
		})
	}
}

func TestManager_Wait(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
		mockFn  func() *Manager
	}{
		{
			name:    "Success",
			wantErr: nil,
			mockFn: func() *Manager {
				return NewManager(2)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := tt.mockFn()

			err := g.Wait()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
