package encode

import (
	"github.com/inovacc/base58"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	str := "hello world"
	enc := base58.StdEncoding.EncodeToString([]byte(str))

	if enc != "StV1DL6CwTryKyV" {
		t.Errorf("Expected %s, got %s", "StV1DL6CwTryKyV", enc)
	}
}

func TestBase58Decode(t *testing.T) {
	data := []byte("Hello, World!")
	encoded := base58.StdEncoding.EncodeToString(data)
	decoded, err := base58.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Errorf("Error decoding string: %v", err)
		return
	}

	if string(decoded) != string(data) {
		t.Errorf("Decoded string does not match original: got %s, want %s", decoded, data)
	}
}
