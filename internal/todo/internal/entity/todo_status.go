package entity

import (
	"database/sql/driver"
	"errors"
)

var ErrScanTodoStatus = errors.New("failed to scan todo status")

type TodoStatus int

const (
	TodoStatusUnknown TodoStatus = iota
	TodoStatusInitiate
	TodoStatusInprogres
	TodoStatusDrop
	TodoStatusDone
)

func ParseTodoStatus(s string) TodoStatus {
	switch s {
	case TodoStatusInitiate.String():
		return TodoStatusInitiate
	case TodoStatusInprogres.String():
		return TodoStatusInprogres
	case TodoStatusDrop.String():
		return TodoStatusDrop
	case TodoStatusDone.String():
		return TodoStatusDone
	default:
		return TodoStatusUnknown
	}
}

func (ts TodoStatus) String() string {
	return [...]string{
		"UNKNOWN",
		"INITIATE",
		"IN_PROGRESS",
		"DROP",
		"DONE",
	}[ts]
}

func (ts *TodoStatus) Scan(value any) error {
	col, ok := value.([]byte)
	if !ok {
		return ErrScanTodoStatus
	}

	*ts = ParseTodoStatus(string(col))

	return nil
}

func (ts TodoStatus) Value() (driver.Value, error) {
	return ts.String(), nil
}
