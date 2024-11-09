package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyVal(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want Field
	}{
		{
			name: "SuccessString",
			args: args{key: "key1", value: "string"},
			want: Field{key: "key1", value: "string"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := KeyVal(tt.args.key, tt.args.value)
			assert.Equal(t, tt.want, got)
		})
	}
}
