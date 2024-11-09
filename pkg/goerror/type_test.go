package goerror

import "testing"

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		et   Type
		want string
	}{
		{
			name: "Validation",
			et:   TypeValidation,
			want: "ERROR_TYPE_VALIDATION",
		},
		{
			name: "Business",
			et:   TypeBusiness,
			want: "ERROR_TYPE_BUSINESS",
		},
		{
			name: "Server",
			et:   TypeServer,
			want: "ERROR_TYPE_SERVER",
		},
		{
			name: "Fallback",
			et:   -1,
			want: "ERROR_TYPE_UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.et.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
