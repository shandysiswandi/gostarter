package domain

import (
	"github.com/shandysiswandi/goreng/enum"
)

type TodoStatus int

const (
	TodoStatusUnknown TodoStatus = iota
	TodoStatusInitiate
	TodoStatusInProgress
	TodoStatusDrop
	TodoStatusDone
)

func (ts TodoStatus) Values() map[enum.Enumerate]string {
	return map[enum.Enumerate]string{
		TodoStatusUnknown:    "UNKNOWN",
		TodoStatusInitiate:   "INITIATE",
		TodoStatusInProgress: "IN_PROGRESS",
		TodoStatusDrop:       "DROP",
		TodoStatusDone:       "DONE",
	}
}
