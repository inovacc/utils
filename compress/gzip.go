package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func GzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func(r *gzip.Reader) {
		_ = r.Close()
	}(r)
	return io.ReadAll(r)
}
