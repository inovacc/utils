package password

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	v, err := HashPassword("test")
	if err != nil {
		t.Error(err)
		return
	}

	if v == "" {
		t.Error("hash is empty")
		return
	}

	if !CheckPasswordHash("test", v) {
		t.Error("hash not match")
		return
	}
}
