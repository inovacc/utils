package file2image

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

const (
	HeaderSize     = 2 + 4 + 6 + 2 // kind (2) + crc (4) + index (6) + dataLen (2)
	ChunkAlignment = 1024
)

type kindOp rune

const (
	KindOperationStart kindOp = iota + 1
	KindOperationFragment
	KindOperationEnd
)

type Header struct {
	Kind    kindOp
	Crc     uint32
	Index   string
	DataLen uint16
}

func (m *Header) Encode(data []byte) []byte {
	m.Crc = crc32.ChecksumIEEE(data)
	m.DataLen = uint16(len(data))

	arr := make([]byte, HeaderSize)
	offset := 0

	binary.BigEndian.PutUint16(arr[offset:], uint16(m.Kind))
	offset += 2

	binary.BigEndian.PutUint32(arr[offset:], m.Crc)
	offset += 4

	idxBytes := make([]byte, 6)
	copy(idxBytes, m.Index)
	copy(arr[offset:], idxBytes)
	offset += 6

	binary.BigEndian.PutUint16(arr[offset:], m.DataLen)
	return arr
}

func (m *Header) Decode(data []byte) error {
	if len(data) < HeaderSize {
		return fmt.Errorf("invalid metadata size")
	}
	offset := 0

	m.Kind = kindOp(binary.BigEndian.Uint16(data[offset:]))
	offset += 2

	m.Crc = binary.BigEndian.Uint32(data[offset:])
	offset += 4

	m.Index = string(data[offset : offset+6])
	offset += 6

	m.DataLen = binary.BigEndian.Uint16(data[offset:])
	return nil
}
