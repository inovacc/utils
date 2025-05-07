package file2image

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const DataSizePerChunk = 800 // bytes of data per chunk (must be < 1024 - HeaderSize)

func TestNewChunk(t *testing.T) {
	filename := "testdata/test.jpg"
	chunker := NewChunker(filename, DataSizePerChunk)

	chunks, err := chunker.Chunks()
	if err != nil {
		log.Fatal(err)
	}

	for chunk := range chunks {
		meta, data, err := chunk.Parse()
		if err != nil {
			log.Fatalf("invalid chunk: %v", err)
		}
		fmt.Printf("Chunk [%s] Kind=%d Size=%d CRC=%08X\n", meta.Index, meta.Kind, len(data), meta.Crc)

		os.WriteFile(fmt.Sprintf("testdata/chunk_%s.bin", meta.Index), data, 0644)
	}
}
