package chunk

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
)

const (
	max48bit   = 0xFFFFFFFFFFFF    // 281,474,976,710,655
	headerSize = 2 + 4 + 6 + 8 + 2 // kind(2) + crc(4) + index(6) + id(8) + dataLen(2)
	Alignment  = 1024
)

func HeaderSize() int {
	return headerSize
}

type kindOp rune

const (
	KindOperationStart kindOp = iota + 1
	KindOperationFragment
	KindOperationEnd
)

func (k kindOp) String() string {
	switch k {
	case KindOperationStart:
		return "START"
	case KindOperationFragment:
		return "FRAGMENT"
	case KindOperationEnd:
		return "END"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(k))
	}
}

func (k kindOp) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *kindOp) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "START":
		*k = KindOperationStart
	case "FRAGMENT":
		*k = KindOperationFragment
	case "END":
		*k = KindOperationEnd
	default:
		return fmt.Errorf("invalid kind: %q", s)
	}
	return nil
}

type Header struct {
	ID     uint64 `json:"id"`
	Kind   kindOp `json:"kind"`
	Crc    uint32 `json:"crc"`
	Index  uint64 `json:"index"`
	Length uint16 `json:"length"`
}

func NewHeader(id uint64, kind kindOp, index uint64) (*Header, error) {
	if index > max48bit {
		return nil, fmt.Errorf("index %d exceeds 48-bit limit", index)
	}

	return &Header{
		ID:    id,
		Kind:  kind,
		Crc:   0,
		Index: index,
	}, nil
}

func (h *Header) ToJSON() (string, error) {
	data, err := json.MarshalIndent(h, "", "  ")
	return string(data), err
}

func (h *Header) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), h)
}

func (h *Header) Encode(data []byte) ([]byte, error) {
	if h.Index > max48bit {
		return nil, fmt.Errorf("index %d exceeds 48-bit limit (max %d)", h.Index, max48bit)
	}

	h.Crc = crc32.ChecksumIEEE(data)
	h.Length = uint16(len(data))

	buf := new(bytes.Buffer)
	buf.Grow(HeaderSize())

	if err := binary.Write(buf, binary.BigEndian, uint16(h.Kind)); err != nil {
		return nil, fmt.Errorf("failed to encode kind: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, h.Crc); err != nil {
		return nil, fmt.Errorf("failed to encode crc: %w", err)
	}

	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes[:], h.Index)
	if _, err := buf.Write(idxBytes[2:]); err != nil {
		return nil, fmt.Errorf("failed to encode index: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, h.ID); err != nil {
		return nil, fmt.Errorf("failed to encode id: %w", err)
	}

	if err := binary.Write(buf, binary.BigEndian, h.Length); err != nil {
		return nil, fmt.Errorf("failed to encode dataLen: %w", err)
	}

	return buf.Bytes(), nil
}

func (h *Header) Decode(data []byte) error {
	if len(data) < HeaderSize() {
		return fmt.Errorf("invalid header size: got %d, want at least %d", len(data), HeaderSize())
	}

	buf := bytes.NewReader(data)

	var kind uint16
	if err := binary.Read(buf, binary.BigEndian, &kind); err != nil {
		return fmt.Errorf("failed to decode kind: %w", err)
	}
	h.Kind = kindOp(kind)

	if err := binary.Read(buf, binary.BigEndian, &h.Crc); err != nil {
		return fmt.Errorf("failed to decode crc: %w", err)
	}

	var compactIdx [6]byte
	if _, err := io.ReadFull(buf, compactIdx[:]); err != nil {
		return fmt.Errorf("failed to decode index: %w", err)
	}

	var fullIdx [8]byte
	copy(fullIdx[2:], compactIdx[:])
	h.Index = binary.BigEndian.Uint64(fullIdx[:])

	if h.Index > max48bit {
		return fmt.Errorf("decoded index %d exceeds 48-bit range", h.Index)
	}

	if err := binary.Read(buf, binary.BigEndian, &h.ID); err != nil {
		return fmt.Errorf("failed to decode id: %w", err)
	}

	if err := binary.Read(buf, binary.BigEndian, &h.Length); err != nil {
		return fmt.Errorf("failed to decode length: %w", err)
	}

	return nil
}
