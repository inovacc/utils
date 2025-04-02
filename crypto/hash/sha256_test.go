package hash

import "testing"

func TestHashSHA256(t *testing.T) {
	h := Sha256("test")
	if h != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Errorf("Expected 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08, got %s", h)
		return
	}

	h = Sha256Bytes([]byte("test"))
	if h != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Errorf("Expected 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08, got %s", h)
		return
	}
}
