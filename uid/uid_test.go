package uid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/inovacc/ksuid"
)

func TestGenerateUUID(t *testing.T) {
	v := GenerateUUID()
	if len(v) != 36 {
		t.Errorf("Expected UUID length of 36, got %d", len(v))
		return
	}

	if _, err := uuid.Parse(v); err != nil {
		t.Errorf("Expected valid UUID, got error: %v", err)
		return
	}
}

func TestGenerateKSUID(t *testing.T) {
	v := GenerateKSUID()
	if len(v) != 27 {
		t.Errorf("Expected KSUID length of 27, got %d", len(v))
		return
	}

	if _, err := ksuid.Parse(v); err != nil {
		t.Errorf("Expected valid KSUID, got error: %v", err)
		return
	}
}
