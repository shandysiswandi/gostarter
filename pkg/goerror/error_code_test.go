package goerror

import "testing"

func TestCode_String(t *testing.T) {
	tests := []struct {
		name string
		c    Code
		want string
	}{
		{
			name: "CodeUnknown",
			c:    CodeUnknown,
			want: "ERROR_CODE_UNKNOWN",
		},
		{
			name: "CodeInvalidFormat",
			c:    CodeInvalidFormat,
			want: "ERROR_CODE_INVALID_FORMAT",
		},
		{
			name: "CodeInvalidInput",
			c:    CodeInvalidInput,
			want: "ERROR_CODE_INVALID_INPUT",
		},
		{
			name: "CodeNotFound",
			c:    CodeNotFound,
			want: "ERROR_CODE_NOT_FOUND",
		},
		{
			name: "CodeConflict",
			c:    CodeConflict,
			want: "ERROR_CODE_CONFLICT",
		},
		{
			name: "CodeUnauthorized",
			c:    CodeUnauthorized,
			want: "ERROR_CODE_UNAUTHORIZED",
		},
		{
			name: "CodeForbidden",
			c:    CodeForbidden,
			want: "ERROR_CODE_FORBIDDEN",
		},
		{
			name: "CodeContentTooLarge",
			c:    CodeContentTooLarge,
			want: "ERROR_CODE_CONTENT_TOO_LARGE",
		},
		{
			name: "CodeTimeout",
			c:    CodeTimeout,
			want: "ERROR_CODE_UNKNOWN",
		},
		{
			name: "CodeInternal",
			c:    CodeInternal,
			want: "ERROR_CODE_INTERNAL",
		},
		{
			name: "Fallback",
			c:    -1,
			want: "ERROR_CODE_UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Code.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
