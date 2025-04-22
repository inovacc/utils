package password

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given plain text password using the bcrypt algorithm.
// It returns the hashed password as a string and any error encountered during hashing.
//
// Example:
//
//	hash, err := HashPassword("mySecret123")
//	if err != nil { ... }
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hashed password with its possible plain text equivalent.
// Returns true if the password matches the hash, false otherwise.
//
// Example:
//
//	match := CheckPasswordHash("mySecret123", storedHash)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
