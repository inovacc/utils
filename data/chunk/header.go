package chunk

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

const (
	headerSize = 2 + 4 + 6 + 8 + 2 // kind(2) + crc(4) + index(6) + id(8) + dataLen(2)
	Alignment  = 1024
)

type kindOp rune

const (
	KindOperationStart kindOp = iota + 1
	KindOperationFragment
	KindOperationEnd
)

type Header struct {
	ID      int
	Kind    kindOp
	Crc     uint32
	Index   string
	DataLen uint16
}

func NewHeader(id int, kind kindOp, index int) *Header {
	return &Header{
		ID:      id,
		Kind:    kind,
		Crc:     0,
		Index:   fmt.Sprintf("%06X", index),
		DataLen: 0,
	}
}

func (m *Header) HeaderSize() int {
	return headerSize
}

func (m *Header) Encode(data []byte) []byte {
	m.Crc = crc32.ChecksumIEEE(data)
	m.DataLen = uint16(len(data))

	arr := make([]byte, m.HeaderSize())
	offset := 0

	binary.BigEndian.PutUint16(arr[offset:], uint16(m.Kind))
	offset += 2

	binary.BigEndian.PutUint32(arr[offset:], m.Crc)
	offset += 4

	idxBytes := make([]byte, 6)
	copy(idxBytes, m.Index)
	copy(arr[offset:], idxBytes)
	offset += 6

	binary.BigEndian.PutUint64(arr[offset:], uint64(m.ID))
	offset += 8

	binary.BigEndian.PutUint16(arr[offset:], m.DataLen)
	return arr
}

func (m *Header) Decode(data []byte) error {
	if len(data) < m.HeaderSize() {
		return fmt.Errorf("invalid metadata size")
	}
	offset := 0

	m.Kind = kindOp(binary.BigEndian.Uint16(data[offset:]))
	offset += 2

	m.Crc = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	m.Index = string(data[offset : offset+6])
	offset += 6

	m.ID = int(binary.BigEndian.Uint64(data[offset:]))
	offset += 8

	m.DataLen = binary.BigEndian.Uint16(data[offset:])
	return nil
}
