package uid

import (
	"github.com/google/uuid"
	"github.com/inovacc/ksuid"
)

// GenerateUUID returns a new RFC 4122 UUID (Universally Unique Identifier) as a string.
// This is useful for creating unique keys, tokens, or identifiers in distributed systems.
//
// Example:
//
//	id := GenerateUUID()
//	fmt.Println(id) // "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
func GenerateUUID() string {
	return uuid.NewString()
}

// GenerateKSUID returns a new KSUID (K-Sortable Unique Identifier) as a string.
// KSUIDs are 27-character, time-sortable identifiers that include a timestamp and random payload.
// This is ideal for systems that benefit from roughly ordered unique IDs (e.g., logs, events).
//
// Example:
//
//	kid := GenerateKSUID()
//	fmt.Println(kid) // "0ujtsYcgvSTl8PAuAdqWYSMnLOv"
func GenerateKSUID() string {
	return ksuid.NewString()
}
