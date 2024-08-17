package gql

import (
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
)

func getString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}

	return ""
}

func getStatusString(status *ql.Status) string {
	if status != nil && status.IsValid() {
		return status.String()
	}

	return ""
}
