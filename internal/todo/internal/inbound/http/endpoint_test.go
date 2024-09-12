package http

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestRegisterRESTEndpoint(t *testing.T) {
	type args struct {
		router *httprouter.Router
		h      *Endpoint
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterRESTEndpoint(tt.args.router, tt.args.h)
		})
	}
}

func TestEndpoint_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Create(tt.args.ctx, tt.args.r)
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
		in1 *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Delete(tt.args.ctx, tt.args.in1)
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
		in1 *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Find(tt.args.ctx, tt.args.in1)
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
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Fetch(tt.args.ctx, tt.args.r)
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
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.UpdateStatus(tt.args.ctx, tt.args.r)
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
		r   *http.Request
	}
	tests := []struct {
		name    string
		e       *Endpoint
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.Update(tt.args.ctx, tt.args.r)
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
