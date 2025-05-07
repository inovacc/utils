package file2image

import (
	"log"
	"testing"
)

func TestNewMetadata(t *testing.T) {
	filename := "testdata/test.jpg"
	metadata, err := NewChunks(filename, 1024)
	if err != nil {
		t.Fatal(err)
	}

	var frames [][]byte
	for obj := range metadata.Data {
		frames = append(frames, obj)
	}

	if err := ShowTxqrSequence(frames); err != nil {
		log.Fatalf("error mostrando secuencia QR: %v", err)
	}

	// if err := metadata.GenerateFramesLive(); err != nil {
	// 	t.Fatal(err)
	// }
}
