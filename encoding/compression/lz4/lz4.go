package lz4

import (
	"bytes"
	"io"

	"github.com/pierrec/lz4/v4"
)

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := lz4.NewWriter(&b)
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
	r := lz4.NewReader(bytes.NewReader(data))
	return io.ReadAll(r)
}
