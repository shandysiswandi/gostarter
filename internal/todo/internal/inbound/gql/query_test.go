package gql

import (
	"context"
	"reflect"
	"testing"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
)

func Test_query_Fetch(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *ql.FetchInput
	}
	tests := []struct {
		name    string
		q       *query
		args    args
		want    []ql.Todo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Fetch(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_query_Find(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		q       *query
		args    args
		want    *ql.Todo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Find(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("query.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("query.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
