package compress

import "github.com/golang/snappy"

func SnappyCompress(data []byte) ([]byte, error) {
	return snappy.Encode(nil, data), nil
}

func SnappyDecompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}
