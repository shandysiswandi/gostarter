package strcase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerSnake(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "lowercase",
			args: args{s: "test"},
			want: "test",
		},
		{
			name: "camelCase",
			args: args{s: "camelCase"},
			want: "camel_case",
		},
		{
			name: "PascalCase",
			args: args{s: "PascalCase"},
			want: "pascal_case",
		},
		{
			name: "multipleUpper",
			args: args{s: "ThisIsATest"},
			want: "this_is_a_test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ToLowerSnake(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
