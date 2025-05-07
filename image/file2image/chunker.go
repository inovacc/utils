package file2image

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
)

type Chunk struct {
	Raw     []byte
	DataLen int
}

// NewChunk creates a new aligned chunk with metadata and data
func NewChunk(kind kindOp, index int, data []byte) Chunk {
	meta := &Header{
		Kind:  kind,
		Index: fmt.Sprintf("%06X", index),
	}

	metaBytes := meta.Encode(data)
	total := len(metaBytes) + len(data)
	remainder := total % ChunkAlignment
	if remainder != 0 {
		total += ChunkAlignment - remainder
	}

	raw := make([]byte, total)
	copy(raw[:HeaderSize], metaBytes)
	copy(raw[HeaderSize:], data)

	return Chunk{
		Raw:     raw,
		DataLen: len(data),
	}
}

// Parse extracts metadata and data from a raw chunk
func (c Chunk) Parse() (*Header, []byte, error) {
	if len(c.Raw) < HeaderSize {
		return nil, nil, fmt.Errorf("chunk too small: got %d bytes", len(c.Raw))
	}

	meta := &Header{}
	if err := meta.Decode(c.Raw); err != nil {
		return nil, nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	// usa meta.DataLen como límite real
	if int(meta.DataLen)+HeaderSize > len(c.Raw) {
		return nil, nil, fmt.Errorf("dataLen out of bounds")
	}

	data := c.Raw[HeaderSize : HeaderSize+int(meta.DataLen)]
	crc := crc32.ChecksumIEEE(data)
	if crc != meta.Crc {
		return nil, nil, fmt.Errorf("crc mismatch: expected %08X, got %08X", meta.Crc, crc)
	}

	return meta, data, nil
}

// Size returns the total size of the chunk
func (c Chunk) Size() int {
	return len(c.Raw)
}

// IsAligned checks if the chunk size is aligned to ChunkAlignment
func (c Chunk) IsAligned() bool {
	return c.Size()%ChunkAlignment == 0
}

type Chunker struct {
	FilePath      string
	ChunkDataSize int
}

func NewChunker(filePath string, chunkDataSize int) *Chunker {
	return &Chunker{
		FilePath:      filePath,
		ChunkDataSize: chunkDataSize,
	}
}

func (c *Chunker) Chunks() (<-chan Chunk, error) {
	file, err := os.Open(c.FilePath)
	if err != nil {
		return nil, err
	}

	out := make(chan Chunk)
	go func() {
		defer func() {
			close(out)
			if err := file.Close(); err != nil {
				log.Printf("error closing file: %v", err)
			}
		}()

		buf := make([]byte, c.ChunkDataSize)
		index := 0

		for {
			n, err := file.Read(buf)
			if n > 0 {
				kind := KindOperationFragment
				if index == 0 {
					kind = KindOperationStart
				}
				chunk := NewChunk(kind, index, buf[:n])
				out <- chunk
				index++
			}
			if err != nil {
				if err == io.EOF {
					endChunk := NewChunk(KindOperationEnd, index, []byte{})
					out <- endChunk
					break
				}
				log.Printf("error reading file: %v\n", err)
				break
			}
		}
	}()

	return out, nil
}
