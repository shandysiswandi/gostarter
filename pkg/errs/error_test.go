package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "no underlying error or message",
			err:  &Error{},
			want: "",
		},
		{
			name: "with underlying error",
			err:  &Error{err: errors.New("underlying error")},
			want: "underlying error",
		},
		{
			name: "with custom message",
			err:  &Error{msg: "custom message"},
			want: "custom message",
		},
		{
			name: "with both underlying error and custom message",
			err:  &Error{err: errors.New("underlying error"), msg: "custom message"},
			want: "underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := &Error{err: underlyingErr}

	if got := err.Unwrap(); !errors.Is(got, underlyingErr) {
		t.Errorf("Unwrap() = %v, want %v", got, underlyingErr)
	}
}

func TestError_Code(t *testing.T) {
	err := &Error{code: CodeInvalidInput}

	if got := err.Code(); got != CodeInvalidInput {
		t.Errorf("Code() = %v, want %v", got, CodeInvalidInput)
	}
}

func TestError_Type(t *testing.T) {
	err := &Error{errType: TypeValidation}

	if got := err.Type(); got != TypeValidation {
		t.Errorf("Type() = %v, want %v", got, TypeValidation)
	}
}

func TestIsValidationError(t *testing.T) {
	validationErr := NewValidation("validation error")
	if !IsValidationError(validationErr) {
		t.Errorf("IsValidationError() = false, want true")
	}

	if IsValidationError(errors.New("")) {
		t.Errorf("IsValidationError() = true, want false")
	}
}

func TestIsBusinessError(t *testing.T) {
	businessErr := NewBusiness("business error")
	if !IsBusinessError(businessErr) {
		t.Errorf("IsBusinessError() = false, want true")
	}

	if IsBusinessError(errors.New("")) {
		t.Errorf("IsBusinessError() = true, want false")
	}
}

func TestIsServerError(t *testing.T) {
	serverErr := NewServer("server error")
	if !IsServerError(serverErr) {
		t.Errorf("IsServerError() = false, want true")
	}

	if IsServerError(errors.New("")) {
		t.Errorf("IsServerError() = true, want false")
	}
}

func TestNewValidation(t *testing.T) {
	err := NewValidation("validation error")
	e, ok := err.(*Error)
	if !ok || e.errType != TypeValidation || e.code != CodeInvalidInput || e.msg != "validation error" {
		t.Errorf("NewValidation() = %v, want &Error{TypeValidation, CodeInvalidInput, 'validation error'}", err)
	}
}

func TestNewBusiness(t *testing.T) {
	err := NewBusiness("business error")
	e, ok := err.(*Error)
	if !ok || e.errType != TypeBusiness || e.code != CodeNotFound || e.msg != "business error" {
		t.Errorf("NewBusiness() = %v, want &Error{TypeBusiness, CodeNotFound, 'business error'}", err)
	}
}

func TestNewServer(t *testing.T) {
	err := NewServer("server error")
	e, ok := err.(*Error)
	if !ok || e.errType != TypeServer || e.code != CodeInternal || e.msg != "server error" {
		t.Errorf("NewServer() = %v, want &Error{TypeServer, CodeInternal, 'server error'}", err)
	}
}

func TestNewValidationFrom(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := NewValidationFrom(underlyingErr)
	e, ok := err.(*Error)
	if !ok || e.err != underlyingErr || e.errType != TypeValidation || e.code != CodeInvalidInput {
		t.Errorf("NewValidationFrom() = %v, want &Error{TypeValidation, CodeInvalidInput, 'underlying error'}", err)
	}
}

func TestNewBusinessFrom(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := NewBusinessFrom(underlyingErr)
	e, ok := err.(*Error)
	if !ok || e.err != underlyingErr || e.errType != TypeBusiness || e.code != CodeNotFound {
		t.Errorf("NewBusinessFrom() = %v, want &Error{TypeBusiness, CodeNotFound, 'underlying error'}", err)
	}
}

func TestNewServerFrom(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := NewServerFrom(underlyingErr)
	e, ok := err.(*Error)
	if !ok || e.err != underlyingErr || e.errType != TypeServer || e.code != CodeInternal {
		t.Errorf("NewServerFrom() = %v, want &Error{TypeServer, CodeInternal, 'underlying error'}", err)
	}
}

func TestWrap(t *testing.T) {
	t.Run("generic_error", func(t *testing.T) {
		underlyingErr := errors.New("internal system")
		err := Wrap(underlyingErr, "errs")
		e, ok := err.(*Error)
		assert.True(t, ok)
		assert.Equal(t, "errs: internal system", e.Error())
	})

	t.Run("custom_error", func(t *testing.T) {
		underlyingErr := New(nil, "internal system", TypeServer, CodeInternal)
		err := Wrap(underlyingErr, "errs")
		e, ok := err.(*Error)
		assert.True(t, ok)
		assert.Equal(t, "errs: internal system", e.Error())
	})
}

func TestNew(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	err := New(underlyingErr, "custom message", TypeBusiness, CodeConflict)
	e, ok := err.(*Error)
	if !ok || e.err != underlyingErr || e.msg != "custom message" || e.errType != TypeBusiness || e.code != CodeConflict {
		t.Errorf("New() = %v, want &Error{custom message, TypeBusiness, CodeConflict, 'underlying error'}", err)
	}
}
