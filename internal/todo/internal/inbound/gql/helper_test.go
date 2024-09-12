package gql

import (
	"testing"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
)

func Test_getString(t *testing.T) {
	type args struct {
		ptr *string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getString(tt.args.ptr); got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStatusString(t *testing.T) {
	type args struct {
		status *ql.Status
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStatusString(tt.args.status); got != tt.want {
				t.Errorf("getStatusString() = %v, want %v", got, tt.want)
			}
		})
	}
}
