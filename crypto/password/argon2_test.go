package password

import (
	"encoding/json"
	"testing"

	"github.com/inovacc/base58"
)

func TestHashPasswordArgon2(t *testing.T) {
	d, err := HashPasswordArgon2("password", nil)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
		return
	}

	var h Hash
	decode, err := base58.StdEncoding.Decode(d)
	if err != nil {
		t.Fatalf("failed to decode hash: %v", err)
		return
	}

	if err := json.Unmarshal(decode, &h); err != nil {
		t.Fatalf("failed to unmarshal hash: %v", err)
		return
	}

	if len(h.Data) != 32 {
		t.Fatalf("expected hash length 32, got %d", len(h.Data))
		return
	}
}

func TestCheckPasswordHashArgon2(t *testing.T) {
	d, err := HashPasswordArgon2("password", nil)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
		return
	}

	ok, err := CheckPasswordHashArgon2(d, "password")
	if err != nil {
		t.Fatalf("failed to check password hash: %v", err)
		return
	}

	if !ok {
		t.Fatal("expected password to match hash")
		return
	}
}
