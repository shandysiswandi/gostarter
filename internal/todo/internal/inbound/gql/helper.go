package gql

import (
	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
)

var errFailedParseToUint = goerror.NewInvalidInput("failed parse id to uint", nil)

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
