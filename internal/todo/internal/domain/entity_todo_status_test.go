package domain

import (
	"testing"

	"github.com/shandysiswandi/goreng/enum"
	"github.com/stretchr/testify/assert"
)

func TestTodoStatus_Values(t *testing.T) {
	tests := []struct {
		name string
		want map[enum.Enumerate]string
	}{
		{
			name: "Success",
			want: map[enum.Enumerate]string{
				TodoStatusUnknown:    "UNKNOWN",
				TodoStatusInitiate:   "INITIATE",
				TodoStatusInProgress: "IN_PROGRESS",
				TodoStatusDrop:       "DROP",
				TodoStatusDone:       "DONE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, TodoStatus(0).Values())
		})
	}
}
