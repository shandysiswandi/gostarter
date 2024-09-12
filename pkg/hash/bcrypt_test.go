package hash

import (
	"testing"
)

func TestBcryptHashVerifier(t *testing.T) {
	hasher := NewBcryptHash(4) // Using a minimum cost factor for bcrypt

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected bool
	}{
		{"ValidHashAndVerify", "super-secret-password", false, true},
		{"InValid", "aVeryLongPasswordWithManyCharacters1234567890aVeryLongPasswordWithManyCha", true, true},
		{"EmptyString", "", false, true},
		{"LongString", "aVeryLongPasswordWithManyCharacters1234567890", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hashed, err := hasher.Hash(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Test Verify with the correct hash
			if got := hasher.Verify(string(hashed), tt.input); got != tt.expected {
				t.Errorf("Verify() = %v, want %v", got, tt.expected)
			}

			// Test Verify with an incorrect password
			if got := hasher.Verify(string(hashed), "incorrect"); got == tt.expected {
				t.Errorf("Verify() = %v, want %v", got, !tt.expected)
			}
		})
	}
}
