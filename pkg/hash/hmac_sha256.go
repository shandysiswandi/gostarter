package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var ErrHMACSHA256Empty = errors.New("input string cannot be empty")

// HMACSHA256Hash implements HMAC SHA-256 hashing with a secret.
type HMACSHA256Hash struct {
	secret string
}

// NewHMACSHA256Hash initializes HMACSHA256Hash with a given secret.
func NewHMACSHA256Hash(secret string) *HMACSHA256Hash {
	return &HMACSHA256Hash{secret: secret}
}

// Hash generates an HMAC SHA-256 hash of the input string.
func (h *HMACSHA256Hash) Hash(str string) ([]byte, error) {
	if str == "" {
		return nil, ErrHMACSHA256Empty
	}

	hh := hmac.New(sha256.New, []byte(h.secret))
	_, err := hh.Write([]byte(str))
	if err != nil {
		return nil, err
	}

	return []byte(hex.EncodeToString(hh.Sum(nil))), nil
}

// Verify checks if a given plaintext string matches the hashed value.
func (h *HMACSHA256Hash) Verify(hashedHex, str string) bool {
	hashed, err := h.Hash(str)
	if err != nil {
		return false
	}

	return hmac.Equal([]byte(hashedHex), []byte(hex.EncodeToString(hashed)))
}
