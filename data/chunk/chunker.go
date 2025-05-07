package chunk

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Chunk struct {
	Raw []byte
}

// NewChunk creates a new aligned chunk with metadata and data
func NewChunk(id int, kind kindOp, index int, data []byte) Chunk {
	meta := NewHeader(id, kind, index)

	metaBytes := meta.Encode(data)
	total := len(metaBytes) + len(data)
	remainder := total % Alignment
	if remainder != 0 {
		total += Alignment - remainder
	}

	raw := make([]byte, total)
	copy(raw[:HeaderSize], metaBytes)
	copy(raw[HeaderSize:], data)

	return Chunk{
		Raw: raw,
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

// IsAligned checks if the chunk size is aligned to Alignment
func (c Chunk) IsAligned() bool {
	return c.Size()%Alignment == 0
}

type Chunker struct {
	ID            int
	FilePath      string
	ChunkDataSize int
}

func NewChunker(filePath string, chunkDataSize int) *Chunker {
	chunker := &Chunker{
		FilePath:      filePath,
		ChunkDataSize: chunkDataSize,
	}
	chunker.ID = chunker.headerID()
	return chunker
}

func (c *Chunker) headerID() int {
	now := time.Now().UnixNano()
	h := fnv.New64a()
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(now))
	h.Write(b)
	return int(h.Sum64())
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
				chunk := NewChunk(c.ID, kind, index, buf[:n])
				out <- chunk
				index++
			}
			if err != nil {
				if err == io.EOF {
					endChunk := NewChunk(c.ID, KindOperationEnd, index, []byte{})
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

func (c *Chunker) RestoreFromChunks(rawChunks [][]byte) error {
	file, err := os.Create(c.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file for restore: %w", err)
	}
	defer file.Close()

	type parsedChunk struct {
		index string
		data  []byte
		kind  kindOp
	}

	var parsed []parsedChunk
	for _, raw := range rawChunks {
		chunk := Chunk{Raw: raw}
		meta, data, err := chunk.Parse()
		if err != nil {
			return fmt.Errorf("failed to parse chunk: %w", err)
		}

		parsed = append(parsed, parsedChunk{
			index: meta.Index,
			data:  data,
			kind:  meta.Kind,
		})
	}

	sort.Slice(parsed, func(i, j int) bool {
		idx1, _ := strconv.ParseInt(parsed[i].index, 16, 32)
		idx2, _ := strconv.ParseInt(parsed[j].index, 16, 32)
		return idx1 < idx2
	})

	for _, chunk := range parsed {
		if chunk.kind != KindOperationEnd {
			if _, err := file.Write(chunk.data); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
		}
	}

	return nil
}

func (c *Chunker) RestoreFromFolder(folderPath string) error {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read folder: %w", err)
	}

	var rawChunks [][]byte
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(folderPath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read chunk file %s: %w", path, err)
		}
		rawChunks = append(rawChunks, data)
	}

	return c.RestoreFromChunks(rawChunks)
}
