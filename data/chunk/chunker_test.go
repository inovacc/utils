package chunk

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/inovacc/utils/v2/io/fileutil"
	"github.com/inovacc/utils/v2/random/random"
)

func TestNewChunk(t *testing.T) {
	chunker := NewChunker("testdata/test.jpg", 4000)

	chunks, err := chunker.Chunks()
	if err != nil {
		log.Fatal(err)
	}

	// var holdChunks [][]byte // temporally hold chunks simulating some kind of transmission

	for chunk := range chunks {
		// holdChunks = append(holdChunks, chunk.Raw)

		meta, data, err := chunk.Parse()
		if err != nil {
			log.Fatalf("invalid chunk: %v", err)
		}
		fmt.Printf("Chunk [%s] Kind=%d Size=%d CRC=%08X\n", meta.Index, meta.Kind, len(data), meta.Crc)

		os.WriteFile(fmt.Sprintf("testdata/chunks/%s.bin", random.RandomString(8)), data, 0644)
	}

	// chunker = NewChunker("testdata/test_reconstructed.jpg", 2000)
	//
	// if err := chunker.RestoreFromChunks(holdChunks); err != nil {
	// 	t.Fatalf("restore failed: %v", err)
	// }

	chunker = NewChunker("testdata/recovered.jpg", 2000)

	if err := chunker.RestoreFromFolder("testdata/chunks"); err != nil {
		log.Fatalf("restore failed: %v", err)
	}
	fmt.Println("✅ Restore from folder complete.")

	// compare files using hash
	if err := fileutil.CompareFiles("testdata/test.jpg", "testdata/recovered.jpg"); err != nil {
		t.Fatalf("files are different: %v", err)
	}
}
