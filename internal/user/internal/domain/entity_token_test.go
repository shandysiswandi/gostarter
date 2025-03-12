package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken_Table(t *testing.T) {
	tests := []struct {
		name string
		t    Token
		want string
	}{
		{
			name: "Success",
			t:    Token{},
			want: "tokens",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.t.Table()
			assert.Equal(t, tt.want, got)
		})
	}
}
