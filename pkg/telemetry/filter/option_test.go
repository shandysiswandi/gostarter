package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithHeaders(t *testing.T) {
	type args struct {
		header []string
	}
	tests := []struct {
		name   string
		args   args
		want   []string
		mockFn func(a args) *Filter
	}{
		{
			name: "Success",
			args: args{header: []string{"test"}},
			want: []string{"test"},
			mockFn: func(a args) *Filter {
				f := &Filter{}
				WithHeaders(a.header...)(f)
				return f
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.mockFn(tt.args).headers
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithQueries(t *testing.T) {
	type args struct {
		query []string
	}
	tests := []struct {
		name   string
		args   args
		want   []string
		mockFn func(a args) *Filter
	}{
		{
			name: "Success",
			args: args{query: []string{"test"}},
			want: []string{"test"},
			mockFn: func(a args) *Filter {
				f := &Filter{}
				WithQueries(a.query...)(f)
				return f
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.mockFn(tt.args).queries
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithFields(t *testing.T) {
	type args struct {
		field []string
	}
	tests := []struct {
		name   string
		args   args
		want   []string
		mockFn func(a args) *Filter
	}{
		{
			name: "Success",
			args: args{field: []string{"test"}},
			want: []string{"test"},
			mockFn: func(a args) *Filter {
				f := &Filter{}
				WithFields(a.field...)(f)
				return f
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.mockFn(tt.args).fields
			assert.Equal(t, tt.want, got)
		})
	}
}
