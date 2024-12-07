package framework

import (
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/graphql"
)

func Test_transportPOST_Supports(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		tr   transportPOST
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Supports(tt.args.r); got != tt.want {
				t.Errorf("transportPOST.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transportPOST_Do(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		exec graphql.GraphExecutor
	}
	tests := []struct {
		name string
		tr   transportPOST
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.Do(tt.args.w, tt.args.r, tt.args.exec)
		})
	}
}
