package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash implements the HashVerifier interface using bcrypt.
type BcryptHash struct {
	// cost specifies the bcrypt cost factor. Higher values are more secure but slower.
	cost int
}

// NewBcryptHash creates a new BcryptHashVerifier with the specified cost factor.
// The cost factor determines the computational complexity of the hashing process.
func NewBcryptHash(cost int) *BcryptHash {
	return &BcryptHash{cost: cost}
}

// Hash hashes the plaintext string using bcrypt and returns the hashed value.
// It returns an error if the hashing process fails.
func (h *BcryptHash) Hash(str string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), h.cost)

	return hashed, err
}

// Verify compares the hashed value with the plaintext string using bcrypt.
// It returns true if the plaintext string matches the hashed value, otherwise false.
func (h *BcryptHash) Verify(hashed, str string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
}
