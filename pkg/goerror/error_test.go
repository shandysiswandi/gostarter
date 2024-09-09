package goerror

import (
	"reflect"
	"testing"

	"google.golang.org/grpc/status"
)

func TestGoError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("GoError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoError_String(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("GoError.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoError_Msg(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Msg(); got != tt.want {
				t.Errorf("GoError.Msg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoError_Type(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want Type
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Type(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoError.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoError_Code(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want Code
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Code(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoError.Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		e       *GoError
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("GoError.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoError_GRPCStatus(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want *status.Status
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GRPCStatus(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoError.GRPCStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.StatusCode(); got != tt.want {
				t.Errorf("GoError.StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		err  error
		msg  string
		et   Type
		code Code
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New(tt.args.err, tt.args.msg, tt.args.et, tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
