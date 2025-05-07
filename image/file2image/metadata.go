package file2image

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

const (
	MetaSize       = 2 + 4 + 6 // minimal metadata in bytes
	ChunkAlignment = 1024
)

type kindOp rune

const (
	KindOperationStart kindOp = iota + 1
	KindOperationFragment
	KindOperationEnd
)

type Metadata struct {
	Kind  kindOp
	Crc   uint32
	Index string
}

func (m *Metadata) EncodeToBytes(data []byte) []byte {
	arr := make([]byte, MetaSize)
	offset := 0

	// kind (2 bytes)
	binary.BigEndian.PutUint16(arr[offset:], uint16(m.Kind))
	offset += 2

	// crc (4 bytes)
	m.Crc = crc32.ChecksumIEEE(data)
	binary.BigEndian.PutUint32(arr[offset:], m.Crc)
	offset += 4

	// index (6 bytes)
	idxBytes := []byte(m.Index)
	if len(idxBytes) > 6 {
		idxBytes = idxBytes[:6]
	}
	copy(arr[offset:], idxBytes)
	return arr
}

func DecodeMetadata(data []byte) (*Metadata, error) {
	if len(data) < MetaSize {
		return nil, fmt.Errorf("invalid metadata size")
	}

	offset := 0
	kind := kindOp(binary.BigEndian.Uint16(data[offset:]))
	offset += 2
	crc := binary.BigEndian.Uint32(data[offset:])
	offset += 4
	index := string(data[offset : offset+6])
	return &Metadata{Kind: kind, Crc: crc, Index: index}, nil
}

type Chunk struct {
	Raw []byte
}

func NewChunk(kind kindOp, index int, data []byte) Chunk {
	meta := &Metadata{
		Kind:  kind,
		Index: fmt.Sprintf("%06X", index),
	}
	metaBytes := meta.EncodeToBytes(data)
	raw := AlignedChunk(metaBytes, data)

	return Chunk{Raw: raw}
}

func AlignedChunk(meta, data []byte) []byte {
	total := len(meta) + len(data)
	remainder := total % ChunkAlignment
	if remainder != 0 {
		total += ChunkAlignment - remainder
	}

	chunk := make([]byte, total)
	copy(chunk[:len(meta)], meta)
	copy(chunk[len(meta):], data)

	return chunk
}

func ParseChunk(chunk []byte) (*Metadata, []byte, error) {
	if len(chunk) < MetaSize {
		return nil, nil, fmt.Errorf("chunk too small: got %d bytes, need at least %d", len(chunk), MetaSize)
	}

	// Parse Metadata
	meta, err := DecodeMetadata(chunk[:MetaSize])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	// Data is whatever remains
	data := chunk[MetaSize:]

	// Optional: Validate CRC
	calculatedCRC := crc32.ChecksumIEEE(data)
	if meta.Crc != calculatedCRC {
		return nil, nil, fmt.Errorf("crc mismatch: expected %08X, got %08X", meta.Crc, calculatedCRC)
	}

	return meta, data, nil
}

func IsValidChunkSize(size int) bool {
	return size >= MetaSize && size%ChunkAlignment == 0
}
