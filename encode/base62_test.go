package encode

import (
	"github.com/inovacc/base62"
	"testing"
)

func TestBase62Encode(t *testing.T) {
	data := []byte("Hello, World!")
	encoded := base62.Encode(data)

	// Check if the encoded string is not empty
	if encoded == "" {
		t.Error("Encoded string is empty")
		return
	}

	// Check if the encoded string is a valid base62 string
	for _, char := range encoded {
		if (char < '0' || char > '9') && (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
			t.Errorf("Encoded string contains invalid character: %c", char)
			return
		}
	}
}

func TestBase62Decode(t *testing.T) {
	data := []byte("Hello, World!")
	encoded := base62.Encode(data)
	decoded, err := base62.Decode(encoded)

	if err != nil {
		t.Errorf("Error decoding string: %v", err)
		return
	}

	if string(decoded) != string(data) {
		t.Errorf("Decoded string does not match original: got %s, want %s", decoded, data)
	}
}
