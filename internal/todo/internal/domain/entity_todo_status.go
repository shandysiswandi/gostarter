package domain

import (
	"database/sql/driver"
	"errors"
	"strings"
)

var ErrScanTodoStatus = errors.New("failed to scan todo status")

type TodoStatus int

const (
	TodoStatusUnknown TodoStatus = iota
	TodoStatusInitiate
	TodoStatusInProgress
	TodoStatusDrop
	TodoStatusDone
)

func ParseTodoStatus(s string) TodoStatus {
	sts := strings.TrimPrefix(s, "STATUS_") // for support grpc enum
	sts = strings.ToUpper(sts)

	switch sts {
	case TodoStatusInitiate.String():
		return TodoStatusInitiate
	case TodoStatusInProgress.String():
		return TodoStatusInProgress
	case TodoStatusDrop.String():
		return TodoStatusDrop
	case TodoStatusDone.String():
		return TodoStatusDone
	default:
		return TodoStatusUnknown
	}
}

func (ts TodoStatus) String() string {
	statuses := [...]string{
		"UNKNOWN",
		"INITIATE",
		"IN_PROGRESS",
		"DROP",
		"DONE",
	}

	if ts < TodoStatusUnknown || int(ts) >= len(statuses) {
		return "UNKNOWN"
	}

	return statuses[ts]
}

func (ts *TodoStatus) Scan(value any) error {
	switch col := value.(type) {
	case []byte:
		*ts = ParseTodoStatus(string(col))
	case string:
		*ts = ParseTodoStatus(col)
	default:
		return ErrScanTodoStatus
	}

	return nil
}

func (ts TodoStatus) Value() (driver.Value, error) {
	return ts.String(), nil
}
