package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Table(t *testing.T) {
	tests := []struct {
		name string
		pr   Account
		want string
	}{
		{
			name: "Success",
			pr:   Account{},
			want: "accounts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.pr.Table()
			assert.Equal(t, tt.want, got)
		})
	}
}
