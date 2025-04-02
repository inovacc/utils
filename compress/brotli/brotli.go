package brotli

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"io"
)

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := brotli.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}
