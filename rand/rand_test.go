package rand

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	v := RandomString(20)
	if v == "" {
		t.Error("string is empty")
		return
	}
	t.Log(v)
}

func TestRandomInt(t *testing.T) {
	v := RandomInt(0, 30)
	if v < 0 {
		t.Error("number is not valid")
		return
	}
}

func TestRandomBytes(t *testing.T) {
	v := RandomBytes(20)
	if v == nil {
		t.Error("bytes are empty")
		return
	}
	t.Log(v)
}
