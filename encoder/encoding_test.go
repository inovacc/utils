package encoder

import (
	"log"
	"testing"
)

func TestBase58Encode(t *testing.T) {
	testEncoding := NewEncoding(Base58)

	str := "hello world"
	encoded, err := testEncoding.EncodeStr(str)
	if err != nil {
		log.Fatal(err)
	}

	if encoded != "StV1DL6CwTryKyV" {
		t.Fatalf("Expected %s, got %s", "StV1DL6CwTryKyV", encoded)
	}
}

func TestBase58Decode(t *testing.T) {
	testEncoding := NewEncoding(Base58)

	data := []byte("Hello, World!")
	encoded, err := testEncoding.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	decoded, err := testEncoding.Decode(encoded)
	if err != nil {
		t.Fatalf("Error decoding string: %v", err)
	}

	if string(decoded) != string(data) {
		t.Fatalf("Decoded string does not match original: got %s, want %s", decoded, data)
	}
}

func TestBase62Encode(t *testing.T) {
	testEncoding := NewEncoding(Base62)

	str := "hello world"
	encoded, err := testEncoding.EncodeStr(str)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the encoded string is not empty
	if encoded != "AAwf93rvy4aWQVw" {
		t.Fatalf("Expected %s, got %s", "AAwf93rvy4aWQVw", encoded)
	}

	// Check if the encoded string is a valid base62 string
	for _, char := range encoded {
		if (char < '0' || char > '9') && (char < 'A' || char > 'Z') && (char < 'a' || char > 'z') {
			t.Fatalf("Encoded string contains invalid character: %c", char)
		}
	}
}

func TestBase62Decode(t *testing.T) {
	testEncoding := NewEncoding(Base62)

	data := []byte("Hello, World!")
	encoded, err := testEncoding.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	decoded, err := testEncoding.Decode(encoded)
	if err != nil {
		t.Fatalf("Error decoding string: %v", err)
	}

	if string(decoded) != string(data) {
		t.Fatalf("Decoded string does not match original: got %s, want %s", decoded, data)
	}
}

func TestBase64Encode(t *testing.T) {
	testEncoding := NewEncoding(Base64)

	str := "hello world"
	encoded, err := testEncoding.EncodeStr(str)
	if err != nil {
		log.Fatal(err)
	}

	if encoded != "aGVsbG8gd29ybGQ=" {
		t.Fatalf("Expected %s, got %s", "aGVsbG8gd29ybGQ=", encoded)
	}
}

func TestBase64Decode(t *testing.T) {
	testEncoding := NewEncoding(Base64)

	data := []byte("Hello, World!")
	encoded, err := testEncoding.Encode(data)
	if err != nil {
		log.Fatal(err)
	}

	decoded, err := testEncoding.Decode(encoded)
	if err != nil {
		t.Fatalf("Error decoding string: %v", err)
	}

	if string(decoded) != string(data) {
		t.Fatalf("Decoded string does not match original: got %s, want %s", decoded, data)
	}
}
