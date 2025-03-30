package encode

import (
	"encoding/base64"
	"testing"
)

func TestBase64Encode(t *testing.T) {
	str := "hello world"
	enc := base64.StdEncoding.EncodeToString([]byte(str))

	if enc != "aGVsbG8gd29ybGQ=" {
		t.Errorf("Expected %s, got %s", "aGVsbG8gd29ybGQ=", enc)
	}
}

func TestBase64Decode(t *testing.T) {
	data := []byte("Hello, World!")
	encoded := base64.StdEncoding.EncodeToString(data)
	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Errorf("Error decoding string: %v", err)
		return
	}

	if string(decoded) != string(data) {
		t.Errorf("Decoded string does not match original: got %s, want %s", decoded, data)
	}
}
