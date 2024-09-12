package gql

import (
	"context"
	"reflect"
	"testing"

	ql "github.com/shandysiswandi/gostarter/api/gen-gql/todo"
)

func Test_mutation_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		in  ql.CreateInput
	}
	tests := []struct {
		name    string
		m       *mutation
		args    args
		want    *ql.Todo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Create(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutation.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutation.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mutation_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		m       *mutation
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Delete(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutation.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("mutation.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mutation_UpdateStatus(t *testing.T) {
	type args struct {
		ctx    context.Context
		id     string
		status ql.Status
	}
	tests := []struct {
		name    string
		m       *mutation
		args    args
		want    *ql.UpdateStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutation.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutation.UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mutation_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
		in  ql.UpdateInput
	}
	tests := []struct {
		name    string
		m       *mutation
		args    args
		want    *ql.UpdateResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Update(tt.args.ctx, tt.args.id, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutation.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutation.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
