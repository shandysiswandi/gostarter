package outbound

import (
	"context"
	"testing"

	"github.com/shandysiswandi/gostarter/internal/shortly/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewMapShort(t *testing.T) {
	tests := []struct {
		name string
		want *MapShort
	}{
		{name: "Success", want: NewMapShort()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewMapShort()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapShort_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Short
		wantErr error
		mockFn  func(a args) *MapShort
	}{
		{
			name:    "NotFound",
			args:    args{ctx: context.TODO(), key: "key"},
			want:    nil,
			wantErr: nil,
			mockFn: func(a args) *MapShort {
				return &MapShort{
					data: make(map[string]domain.Short),
				}
			},
		},
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), key: "key"},
			want:    &domain.Short{Key: "key"},
			wantErr: nil,
			mockFn: func(a args) *MapShort {
				data := make(map[string]domain.Short)
				data["key"] = domain.Short{Key: "key"}

				return &MapShort{
					data: data,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ms := tt.mockFn(tt.args)
			got, err := ms.Get(tt.args.ctx, tt.args.key)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapShort_Set(t *testing.T) {
	type args struct {
		ctx   context.Context
		value domain.Short
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		mockFn  func(a args) *MapShort
	}{
		{
			name:    "Success",
			args:    args{ctx: context.TODO(), value: domain.Short{Key: "key"}},
			wantErr: nil,
			mockFn: func(a args) *MapShort {
				return &MapShort{
					data: make(map[string]domain.Short),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ms := tt.mockFn(tt.args)
			err := ms.Set(tt.args.ctx, tt.args.value)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
