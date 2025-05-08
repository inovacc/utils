package splitter

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/inovacc/utils/v2/random/random"
)

type op uint16

const (
	OperationStart op = iota + 1
	OperationFragment
	OperationEnd
	OperationUnknown
)

const (
	max48bit   = 0xFFFFFFFFFFFF
	headerSize = 2 + 4 + 6 + 8 + 2 // kind(2) + crc(4) + index(6) + id(8) + dataLen(2)
	footerSize = 8 + 8 + 2 + 256   // total(8) + size(8) + time(8) + name(256)
)

type ChunkHeader struct {
	ID     uint64
	Kind   op
	Crc    uint32
	Index  uint64
	Length uint16
}

type ChunkFooter struct {
	Total uint64
	Size  int64
	Time  int64
	Name  string
}

type Splitter interface {
	Split(file *os.File, outDir string, size int) error
	Merge(inDir, outDir string) error
	Calculate(file *os.File, size int) uint64
}

type Impl struct {
	Name     string `json:"name"`
	Filename []byte `json:"filename"`
	Time     int64  `json:"timestamp"`
	Total    uint64 `json:"total"`
	Size     int64  `json:"size"`
	NameLen  uint16 `json:"nameLen"`
}

func NewMetaChunk() Splitter {
	return &Impl{}
}

// Calculate the number of chunks in the file
func (i *Impl) Calculate(file *os.File, size int) uint64 {
	buf := make([]byte, size)
	var index uint64

	for {
		n, err := file.Read(buf)
		if n > 0 {
			index++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
	}

	_, _ = file.Seek(0, 0)
	return index
}

// Split the file into chunks of size
func (i *Impl) Split(file *os.File, outDir string, size int) error {
	if len(file.Name()) > math.MaxUint16 {
		return fmt.Errorf("filename too long")
	}

	buf := make([]byte, size)
	var index uint64

	id := i.headerID()
	for {
		n, err := file.Read(buf)
		if n > 0 {
			kind := OperationFragment
			if index == 0 {
				kind = OperationStart
			}

			data, err := i.insert(buf[:n], index, kind, id)
			if err != nil {
				return err
			}

			if err := i.save(outDir, data); err != nil {
				return err
			}
			index++
		}
		if err != nil {
			if err == io.EOF {
				data, err := i.summary(index, OperationEnd, id, file)
				if err != nil {
					return err
				}

				if err := i.save(outDir, data); err != nil {
					return err
				}
				break
			}
			break
		}
	}
	return nil
}

func (i *Impl) Merge(inDir, outDir string) error {
	entries, err := os.ReadDir(inDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	type parsedChunk struct {
		index uint64
		data  []byte
		kind  op
		id    uint64
	}

	var chunks []parsedChunk
	var footerData []byte

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(inDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil || len(data) < headerSize {
			continue
		}

		buf := bytes.NewReader(data)

		header := &ChunkHeader{}

		if err := binary.Read(buf, binary.BigEndian, &header.Kind); err != nil {
			return fmt.Errorf("failed to decode kind: %w", err)
		}

		if err := binary.Read(buf, binary.BigEndian, &header.Crc); err != nil {
			return fmt.Errorf("failed to decode crc: %w", err)
		}

		var compactIdx [6]byte
		if _, err := io.ReadFull(buf, compactIdx[:]); err != nil {
			return fmt.Errorf("failed to decode index: %w", err)
		}

		var fullIdx [8]byte
		copy(fullIdx[2:], compactIdx[:])
		header.Index = binary.BigEndian.Uint64(fullIdx[:])

		if header.Index > max48bit {
			return fmt.Errorf("decoded index %d exceeds 48-bit range", header.Index)
		}

		if err := binary.Read(buf, binary.BigEndian, &header.ID); err != nil {
			return fmt.Errorf("failed to decode id: %w", err)
		}

		if err := binary.Read(buf, binary.BigEndian, &header.Length); err != nil {
			return fmt.Errorf("failed to decode length: %w", err)
		}

		chunkData := data[headerSize : headerSize+header.Length]

		if header.Kind == OperationEnd {
			footerData = data[headerSize:]
			continue
		}

		chunks = append(chunks, parsedChunk{index: header.Index, data: chunkData, kind: header.Kind})
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].index < chunks[j].index
	})

	footer := &ChunkFooter{}
	footerBuf := bytes.NewReader(footerData)
	binary.Read(footerBuf, binary.BigEndian, &footer.Total)
	binary.Read(footerBuf, binary.BigEndian, &footer.Size)
	binary.Read(footerBuf, binary.BigEndian, &footer.Time)
	nameBuf := make([]byte, 256)
	footerBuf.Read(nameBuf)
	footer.Name = string(bytes.Trim(nameBuf, "\x00"))

	file, err := os.Create(filepath.Join(outDir, footer.Name))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, chunk := range chunks {
		if _, err := file.Write(chunk.data); err != nil {
			return err
		}
	}

	size, _ := file.Stat()

	if size.Size() != footer.Size {
		return fmt.Errorf("size mismatch: expected %d, got %d", footer.Size, size.Size())
	}

	return nil
}

func (i *Impl) headerID() uint64 {
	var entropy [16]byte
	_, _ = rand.Read(entropy[:])

	now := time.Now().UnixNano()
	binary.BigEndian.PutUint64(entropy[8:], uint64(now))

	h := fnv.New64a()
	if _, err := h.Write(entropy[:]); err != nil {
		return 0
	}
	return h.Sum64()
}

func (i *Impl) header(data []byte, idx uint64, kind op, id uint64) (*bytes.Buffer, error) {
	headerBytes := new(bytes.Buffer)
	headerBytes.Grow(headerSize)

	if err := binary.Write(headerBytes, binary.BigEndian, kind); err != nil {
		return nil, fmt.Errorf("failed to encode kind: %w", err)
	}

	if err := binary.Write(headerBytes, binary.BigEndian, crc32.ChecksumIEEE(data)); err != nil {
		return nil, fmt.Errorf("failed to encode crc: %w", err)
	}

	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes[:], idx)
	if _, err := headerBytes.Write(idxBytes[2:]); err != nil {
		return nil, fmt.Errorf("failed to encode index: %w", err)
	}

	if err := binary.Write(headerBytes, binary.BigEndian, id); err != nil {
		return nil, fmt.Errorf("failed to encode id: %w", err)
	}

	if err := binary.Write(headerBytes, binary.BigEndian, uint16(len(data))); err != nil {
		return nil, fmt.Errorf("failed to encode dataLen: %w", err)
	}

	return headerBytes, nil
}

func (i *Impl) insert(data []byte, idx uint64, kind op, id uint64) ([]byte, error) {
	if idx > max48bit {
		return nil, fmt.Errorf("index %d exceeds 48-bit limit (max %d)", idx, max48bit)
	}

	headerBytes, err := i.header(data, idx, kind, id)
	if err != nil {
		return nil, err
	}

	raw := make([]byte, headerBytes.Len()+len(data))
	copy(raw[:headerSize], headerBytes.Bytes())
	copy(raw[headerSize:], data)
	return raw, nil
}

func (i *Impl) save(outDir string, data []byte) error {
	name := filepath.Join(outDir, fmt.Sprintf("%s.dat", random.RandomString(8)))
	if err := os.WriteFile(name, data, 0644); err != nil {
		return err
	}
	return nil
}

func (i *Impl) summary(idx uint64, kind op, id uint64, file *os.File) ([]byte, error) {
	if idx > max48bit {
		return nil, fmt.Errorf("index %d exceeds 48-bit limit (max %d)", idx, max48bit)
	}

	size, _ := file.Stat()
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, idx); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, size.Size()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, time.Now().Unix()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, []byte(filepath.Base(file.Name()))); err != nil {
		return nil, err
	}

	headerBytes, err := i.header(buf.Bytes(), idx, kind, id)
	if err != nil {
		return nil, err
	}

	raw := make([]byte, headerSize+footerSize)
	copy(raw[:headerSize], headerBytes.Bytes())
	copy(raw[headerSize:], buf.Bytes())
	return raw, nil
}
