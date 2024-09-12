package hash

import "testing"

func TestArgon2HashVerifier(t *testing.T) {
	hasher := NewArgon2Hash(1, 1, 1, 8) // Example parameters for test

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		expected bool
	}{
		{"Valid Hash and Verify", "super-secret-password", false, true},
		// {"InValid", "", true, false},
		{"Empty String", "", false, true},
		{"Long Password", "aVeryLongPasswordWithManyCharacters1234567890", false, true},
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
