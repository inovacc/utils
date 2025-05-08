package splitter

import (
	"os"
	"testing"
)

func TestNewMetaChunk(t *testing.T) {
	file, err := os.Open("testdata/ubuntu-25.04-desktop-amd64.iso")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	chunk := NewMetaChunk()
	// t.Log(chunk.Calculate(file, 20000))

	t.Log(chunk.Split(file, "testdata/chunks", 2000000))

	chunk.Merge("testdata/chunks", "testdata/output")
}
