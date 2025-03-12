package framework

import (
	"context"
	"reflect"
	"testing"

	"github.com/shandysiswandi/goreng/validation"
	"google.golang.org/grpc"
)

func TestUnaryServerRecovery(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  any
		in2  *grpc.UnaryServerInfo
		next grpc.UnaryHandler
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnaryServerRecovery(tt.args.ctx, tt.args.req, tt.args.in2, tt.args.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnaryServerRecovery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnaryServerRecovery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnaryServerError(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  any
		in2  *grpc.UnaryServerInfo
		next grpc.UnaryHandler
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnaryServerError(tt.args.ctx, tt.args.req, tt.args.in2, tt.args.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnaryServerError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnaryServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnaryServerJWT(t *testing.T) {
	type args struct {
		audience    string
		skipMethods []string
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryServerInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnaryServerJWT(tt.args.audience, tt.args.skipMethods...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnaryServerJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doUnaryServerJWT(t *testing.T) {
	type args struct {
		ctx      context.Context
		req      any
		next     grpc.UnaryHandler
		audience string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doUnaryServerJWT(tt.args.ctx, tt.args.req, tt.args.next, tt.args.audience)
			if (err != nil) != tt.wantErr {
				t.Errorf("doUnaryServerJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doUnaryServerJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnaryServerProtoValidate(t *testing.T) {
	type args struct {
		validator validation.Validator
	}
	tests := []struct {
		name string
		args args
		want grpc.UnaryServerInterceptor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnaryServerProtoValidate(tt.args.validator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnaryServerProtoValidate() = %v, want %v", got, tt.want)
			}
		})
	}
}
