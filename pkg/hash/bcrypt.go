package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHashVerifier implements the HashVerifier interface using bcrypt.
type BcryptHashVerifier struct {
	// Cost specifies the bcrypt cost factor. Higher values are more secure but slower.
	Cost int
}

// NewBcryptHashVerifier creates a new BcryptHashVerifier with the specified cost factor.
// The cost factor determines the computational complexity of the hashing process.
func NewBcryptHashVerifier(cost int) *BcryptHashVerifier {
	return &BcryptHashVerifier{Cost: cost}
}

// Hash hashes the plaintext string using bcrypt and returns the hashed value.
// It returns an error if the hashing process fails.
func (h *BcryptHashVerifier) Hash(str string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), h.Cost)

	return hashed, err
}

// Verify compares the hashed value with the plaintext string using bcrypt.
// It returns true if the plaintext string matches the hashed value, otherwise false.
func (h *BcryptHashVerifier) Verify(hashed, str string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
}
