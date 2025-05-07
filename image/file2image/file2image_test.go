package file2image

import "testing"

func TestNewMetadata(t *testing.T) {
	filename := "testdata/test.jpg"
	metadata, err := NewChunks(filename, 5)
	if err != nil {
		t.Fatal(err)
	}

	if err := metadata.GenerateFramesLive(); err != nil {
		t.Fatal(err)
	}
}
