package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256 returns the SHA-256 hexadecimal hash of the given string input.
// It is useful for generating fixed-length signatures or identifiers from plaintext.
func Sha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Sha256Bytes returns the SHA-256 hexadecimal hash of the given byte slice.
// This is suitable when working with binary data directly instead of strings.
func Sha256Bytes(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
