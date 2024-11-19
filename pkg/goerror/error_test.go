package goerror

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGoError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want string
	}{
		{
			name: "Err",
			e:    &GoError{err: assert.AnError},
			want: "assert.AnError general error for testing",
		},
		{
			name: "Msg",
			e:    &GoError{msg: "error"},
			want: "error",
		},
		{
			name: "TypeValidation",
			e:    &GoError{errType: TypeValidation},
			want: "Validation violation",
		},
		{
			name: "TypeBusiness",
			e:    &GoError{errType: TypeBusiness},
			want: "Logical business not meet with requirement",
		},
		{
			name: "TypeServer",
			e:    &GoError{errType: TypeServer},
			want: "Internal error",
		},
		{
			name: "Unknown",
			e:    &GoError{errType: TypeValidation - 1},
			want: "Unknown error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoError_String(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want string
	}{
		{
			name: "String",
			e:    &GoError{err: assert.AnError},
			want: "Error Type: ERROR_TYPE_VALIDATION, Code: ERROR_CODE_UNKNOWN, Message: , Underlying Error: assert.AnError general error for testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoError_Msg(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want string
	}{
		{
			name: "Msg",
			e:    &GoError{err: assert.AnError},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.Msg()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoError_Type(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want Type
	}{
		{
			name: "Type",
			e:    &GoError{err: assert.AnError},
			want: TypeValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.Type()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoError_Code(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want Code
	}{
		{
			name: "Code",
			e:    &GoError{err: assert.AnError},
			want: CodeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.Code()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		e       *GoError
		wantErr error
	}{
		{
			name:    "Unwrap",
			e:       &GoError{err: assert.AnError},
			wantErr: assert.AnError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.e.Unwrap()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGoError_GRPCStatus(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want *status.Status
	}{
		{
			name: "Unknown",
			e:    &GoError{err: assert.AnError},
			want: status.New(codes.Unknown, ""),
		},
		{
			name: "FailedPrecondition",
			e:    &GoError{code: CodeInvalidInput},
			want: status.New(codes.FailedPrecondition, ""),
		},
		{
			name: "InvalidArgument",
			e:    &GoError{code: CodeInvalidFormat},
			want: status.New(codes.InvalidArgument, ""),
		},
		{
			name: "NotFound",
			e:    &GoError{code: CodeNotFound},
			want: status.New(codes.NotFound, ""),
		},
		{
			name: "Unauthenticated",
			e:    &GoError{code: CodeUnauthorized},
			want: status.New(codes.Unauthenticated, ""),
		},
		{
			name: "PermissionDenied",
			e:    &GoError{code: CodeForbidden},
			want: status.New(codes.PermissionDenied, ""),
		},
		{
			name: "DeadlineExceeded",
			e:    &GoError{code: CodeTimeout},
			want: status.New(codes.DeadlineExceeded, ""),
		},
		{
			name: "AlreadyExists",
			e:    &GoError{code: CodeConflict},
			want: status.New(codes.AlreadyExists, ""),
		},
		{
			name: "Internal",
			e:    &GoError{code: CodeInternal},
			want: status.New(codes.Internal, ""),
		},
		{
			name: "Default",
			e:    &GoError{code: CodeInternal + 1},
			want: status.New(codes.Unknown, ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.GRPCStatus()
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.Code().String(), got.Code().String())
			assert.Equal(t, tt.want.Err().Error(), got.Err().Error())
		})
	}
}

func TestGoError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		e    *GoError
		want int
	}{
		{
			name: "StatusInternalServerError",
			e:    &GoError{err: assert.AnError},
			want: http.StatusInternalServerError,
		},
		{
			name: "StatusBadRequest",
			e:    &GoError{code: CodeInvalidFormat},
			want: http.StatusBadRequest,
		},
		{
			name: "StatusUnprocessableEntity",
			e:    &GoError{code: CodeInvalidInput},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "StatusNotFound",
			e:    &GoError{code: CodeNotFound},
			want: http.StatusNotFound,
		},
		{
			name: "StatusUnauthorized",
			e:    &GoError{code: CodeUnauthorized},
			want: http.StatusUnauthorized,
		},
		{
			name: "StatusForbidden",
			e:    &GoError{code: CodeForbidden},
			want: http.StatusForbidden,
		},
		{
			name: "StatusRequestTimeout",
			e:    &GoError{code: CodeTimeout},
			want: http.StatusRequestTimeout,
		},
		{
			name: "StatusConflict",
			e:    &GoError{code: CodeConflict},
			want: http.StatusConflict,
		},
		{
			name: "Default",
			e:    &GoError{code: CodeUnknown - 1},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.e.StatusCode()
			assert.Equal(t, tt.want, got)
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
		wantErr error
	}{
		{
			name: "New",
			args: args{
				err:  assert.AnError,
				msg:  "general",
				et:   TypeServer,
				code: CodeInternal,
			},
			wantErr: &GoError{
				err:     assert.AnError,
				msg:     "general",
				errType: TypeServer,
				code:    CodeInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := New(tt.args.err, tt.args.msg, tt.args.et, tt.args.code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewServer(t *testing.T) {
	type args struct {
		msg string
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "NewServer",
			args: args{
				msg: "internal",
				err: assert.AnError,
			},
			wantErr: &GoError{
				err:     assert.AnError,
				msg:     "internal",
				errType: TypeServer,
				code:    CodeInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := NewServer(tt.args.msg, tt.args.err)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewServerInternal(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "NewServerInternal",
			args: args{
				err: assert.AnError,
			},
			wantErr: &GoError{
				err:     assert.AnError,
				msg:     "internal server error",
				errType: TypeServer,
				code:    CodeInternal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := NewServerInternal(tt.args.err)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewBusiness(t *testing.T) {
	type args struct {
		msg  string
		code Code
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "NewBusiness",
			args: args{
				msg:  "notfound",
				code: CodeNotFound,
			},
			wantErr: &GoError{
				msg:     "notfound",
				errType: TypeBusiness,
				code:    CodeNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := NewBusiness(tt.args.msg, tt.args.code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewInvalidInput(t *testing.T) {
	type args struct {
		msg string
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "NewInvalidInput",
			args: args{
				msg: "input",
				err: assert.AnError,
			},
			wantErr: &GoError{
				msg:     "input",
				err:     assert.AnError,
				errType: TypeValidation,
				code:    CodeInvalidInput,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := NewInvalidInput(tt.args.msg, tt.args.err)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestNewInvalidFormat(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "NewInvalidFormat",
			args: args{
				msg: "format",
			},
			wantErr: &GoError{
				msg:     "format",
				err:     nil,
				errType: TypeValidation,
				code:    CodeInvalidFormat,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := NewInvalidFormat(tt.args.msg)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
