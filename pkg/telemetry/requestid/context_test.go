package requestid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	type args struct {
		ctx context.Context
		rID string
	}
	tests := []struct {
		name string
		args args
		want func(a args) context.Context
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				rID: "123",
			},
			want: func(a args) context.Context {
				return context.WithValue(a.ctx, requestIDKey{}, a.rID)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Set(tt.args.ctx, tt.args.rID)
			assert.Equal(t, tt.want(tt.args), got)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		ctx  func() context.Context
		want string
	}{
		{
			name: "Success",
			ctx: func() context.Context {
				return context.WithValue(context.Background(), requestIDKey{}, "123")
			},
			want: "123",
		},
		{
			name: "Empty",
			ctx: func() context.Context {
				return context.Background()
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Get(tt.ctx())
			assert.Equal(t, tt.want, got)
		})
	}
}
