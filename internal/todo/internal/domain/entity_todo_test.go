package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodo_ScanColumn(t *testing.T) {
	td := &Todo{}
	tests := []struct {
		name string
		tr   *Todo
		want []any
	}{
		{
			name: "Success",
			tr:   td,
			want: []any{
				&td.ID,
				&td.UserID,
				&td.Title,
				&td.Description,
				&td.Status,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.tr.ScanColumn())
		})
	}
}
