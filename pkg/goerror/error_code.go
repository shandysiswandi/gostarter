package goerror

// Code represents a specific error code for identifying the exact error.
type Code int

// Constants for common error codes. Extend these based on application needs.
const (
	CodeUnknown         Code = iota // Unknown or unspecified error.
	CodeInvalidFormat               // Error code for invalid format.
	CodeInvalidInput                // Error code for invalid input.
	CodeNotFound                    // Error code for resource not found.
	CodeConflict                    // Error code for conflict situations (e.g., duplicate entries).
	CodeUnauthorized                // Error code for unauthorized access.
	CodeForbidden                   // Error code for forbidden actions.
	CodeContentTooLarge             // Error code for content too large.
	CodeTimeout                     // Error code for operation timeout.
	CodeInternal                    // Error code for internal server errors.
)

// String returns the string representation of the error Code.
func (c Code) String() string {
	switch c {
	case CodeInvalidFormat:
		return "ERROR_CODE_INVALID_FORMAT"
	case CodeInvalidInput:
		return "ERROR_CODE_INVALID_INPUT"
	case CodeNotFound:
		return "ERROR_CODE_NOT_FOUND"
	case CodeConflict:
		return "ERROR_CODE_CONFLICT"
	case CodeUnauthorized:
		return "ERROR_CODE_UNAUTHORIZED"
	case CodeForbidden:
		return "ERROR_CODE_FORBIDDEN"
	case CodeContentTooLarge:
		return "ERROR_CODE_CONTENT_TOO_LARGE"
	case CodeUnknown:
		return "ERROR_CODE_UNKNOWN"
	case CodeInternal:
		return "ERROR_CODE_INTERNAL"
	default:
		return "ERROR_CODE_UNKNOWN"
	}
}
