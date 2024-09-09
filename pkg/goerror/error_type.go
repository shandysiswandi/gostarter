package goerror

// Type represents the domain or category of an error.
type Type int

// Constants for different error types.
const (
	TypeValidation Type = iota // Validation errors (e.g., input validation failures).
	TypeBusiness               // Business logic errors (e.g., domain rule violations).
	TypeServer                 // Server-side errors (e.g., database or network issues).
)

// String returns the string representation of the error Type.
func (t Type) String() string {
	switch t {
	case TypeValidation:
		return "ERROR_TYPE_VALIDATION"
	case TypeBusiness:
		return "ERROR_TYPE_BUSINESS"
	case TypeServer:
		return "ERROR_TYPE_SERVER"
	default:
		return "ERROR_TYPE_UNKNOWN"
	}
}
