package chunk

import (
	"hash/crc32"
	"testing"
)

func TestHeaderEncodeDecode(t *testing.T) {
	newChunker := NewChunker("", 0)
	original, _ := NewHeader(newChunker.ID, KindOperationFragment, 7)

	payload := []byte("hello world - this is test data")

	encoded, err := original.Encode(payload)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	if len(encoded) != HeaderSize() {
		t.Errorf("unexpected encoded size: got %d, want %d", len(encoded), HeaderSize())
	}

	var decoded Header
	if err := decoded.Decode(encoded); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	// Validate
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %v, want %v", decoded.ID, original.ID)
	}

	if decoded.Kind != original.Kind {
		t.Errorf("Kind mismatch: got %v, want %v", decoded.Kind, original.Kind)
	}

	if decoded.Index != original.Index {
		t.Errorf("Index mismatch: got %d, want %d", decoded.Index, original.Index)
	}

	if decoded.Length != uint16(len(payload)) {
		t.Errorf("Length mismatch: got %d, want %d", decoded.Length, len(payload))
	}

	if decoded.Crc != original.Crc {
		expectedCRC := crc32.ChecksumIEEE(payload)
		if decoded.Crc != expectedCRC {
			t.Errorf("CRC mismatch: got %08X, want %08X", decoded.Crc, expectedCRC)
		}
	}
}
