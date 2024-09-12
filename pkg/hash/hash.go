/*
Package hash provides interfaces and implementations for hashing and verifying strings.

The package defines the HashVerifier interface, which includes methods for hashing plaintext strings
and verifying hashed values. It includes implementations using different algorithms:

  - BcryptHash: Uses the bcrypt hashing algorithm for secure password hashing.

  - Argon2Hash: Uses the Argon2 key derivation function for secure password hashing.

Example usage:

	// Create a new BcryptHash with a cost factor of 12
	bcrypt := hash.NewBcryptHash(12)

	// Hash a password
	hashedPassword, err := bcrypt.Hash("super-secret-password")
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Printf("Hashed password: %s\n", hashedPassword)

	// Verify the hashed password
	isMatch := bcrypt.Verify(string(hashedPassword), "super-secret-password")
	if isMatch {
	    fmt.Println("Password matches")
	} else {
	    fmt.Println("Password does not match")
	}

	// Create a new Argon2Hash with specific parameters
	argon2 := hash.NewArgon2Hash(1, 64*1024, 4, 32)

	// Hash a password
	hashedPassword, err := argon2.Hash("super-secret-password")
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Printf("Hashed password: %s\n", hashedPassword)

	// Verify the hashed password
	isMatch := argon2.Verify(string(hashedPassword), "super-secret-password")
	if isMatch {
	    fmt.Println("Password matches")
	} else {
	    fmt.Println("Password does not match")
	}
*/
package hash

// HashVerifier defines methods for hashing and verifying strings.
type Hash interface {
	// Hash takes a plaintext string and returns its hashed representation.
	// It may return an error if the hashing process fails.
	Hash(str string) ([]byte, error)

	// Verify checks if the given plaintext string matches the hashed value.
	// It returns true if the plaintext matches the hash, otherwise false.
	Verify(hashed, str string) bool
}
