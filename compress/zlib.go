package compress

import (
	"bytes"
	"compress/zlib"
	"io"
)

func ZlibCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func ZlibDecompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(r)
	return io.ReadAll(r)
}
