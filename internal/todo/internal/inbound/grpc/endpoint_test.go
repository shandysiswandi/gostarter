package grpc

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/shandysiswandi/gostarter/api/gen-proto/todo"
	"github.com/shandysiswandi/gostarter/internal/todo/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/validation"
)

func TestNewEndpoint(t *testing.T) {
	type args struct {
		validator      validation.Validator
		FindUC         domain.Find
		FetchUC        domain.Fetch
		createUC       domain.Create
		deleteUC       domain.Delete
		updateUC       domain.Update
		updateStatusUC domain.UpdateStatus
	}
	tests := []struct {
		name string
		args args
		want *Endpoint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEndpoint(tt.args.validator, tt.args.FindUC, tt.args.FetchUC, tt.args.createUC, tt.args.deleteUC, tt.args.updateUC, tt.args.updateStatusUC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    *pb.CreateResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.DeleteRequest
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    *pb.DeleteResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Find(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.FindRequest
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    *pb.FindResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Find(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Fetch(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.FetchRequest
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    *pb.FetchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Fetch(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_UpdateStatus(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateStatusRequest
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    *pb.UpdateStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.UpdateStatus(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.UpdateRequest
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    *pb.UpdateResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Endpoint.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
