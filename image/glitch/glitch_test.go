package glitch

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/inovacc/utils/v2/crypto/hashing"
)

func TestGlitchEncodeDecode(t *testing.T) {
	g := NewGlitch()

	// Setup test directories and files
	testInput := "/home/dyam/Downloads/gocv-0.41.0.zip"
	testOutputDir := "testdata/frames"
	testReconstructed := "testdata/reconstructed.txt"

	defer os.RemoveAll(testOutputDir)
	defer os.Remove(testReconstructed)

	// Encode a file to images
	if err := g.BlobToImagesFromReader(testInput, testOutputDir); err != nil {
		t.Fatalf("BlobToImages failed: %v", err)
	}

	// Decode images back to file
	pattern := filepath.Join(testOutputDir, "frame_*.png")
	if err := g.ImagesToBlob(pattern, testReconstructed); err != nil {
		t.Fatalf("ImagesToBlob failed: %v", err)
	}

	// Validate content
	got, err := os.ReadFile(testReconstructed)
	if err != nil {
		t.Fatalf("Failed to read reconstructed file: %v", err)
	}

	newHasher := hashing.NewHasher(hashing.SHA256)
	gotStr := newHasher.HashBytes(got)

	source, err := os.ReadFile(testInput)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}

	sourceStr := newHasher.HashBytes(source)

	if gotStr != sourceStr {
		t.Fatal("not equal")
	}
}
