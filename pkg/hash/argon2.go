package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Hash implements the HashVerifier interface using Argon2.
// Argon2 is a modern and secure key derivation function.
type Argon2Hash struct {
	// time is the number of iterations for the Argon2 hashing algorithm.
	time uint32

	// memory is the memory size in KB used by the Argon2 algorithm.
	memory uint32

	// threads is the number of threads used for hashing.
	threads uint8

	// keyLen is the length of the generated key in bytes.
	keyLen uint32
}

// NewArgon2Hash creates a new Argon2HashVerifier with the specified parameters.
// It configures the Argon2 algorithm with the provided time, memory, threads, and key length.
func NewArgon2Hash(time, memory uint32, threads uint8, keyLen uint32) *Argon2Hash {
	return &Argon2Hash{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

// Hash hashes the plaintext string using Argon2 and returns the hashed value.
// It uses a random salt and encodes the result in a format that includes the salt and hash.
func (h *Argon2Hash) Hash(str string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	hash := argon2.IDKey([]byte(str), salt, h.time, h.memory, h.threads, h.keyLen)

	return []byte(hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash)), nil
}

// Verify compares the hashed value with the plaintext string using Argon2.
// It extracts the salt from the hashed value and verifies if the plaintext string matches the hash.
func (h *Argon2Hash) Verify(hashed, str string) bool {
	parts := strings.Split(hashed, ":")
	if len(parts) != 2 {
		return false
	}

	salt, hashHex := parts[0], parts[1]

	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return false
	}

	hashBytes, err := hex.DecodeString(hashHex)
	if err != nil {
		return false
	}

	hash := argon2.IDKey([]byte(str), saltBytes, h.time, h.memory, h.threads, uint32(len(hashBytes)))

	return subtle.ConstantTimeCompare(hash, hashBytes) == 1
}
