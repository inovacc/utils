package zstd

import (
	"bytes"
	"github.com/inovacc/utils/v2/encode"
	"testing"
)

func TestCompress(t *testing.T) {
	data, err := Compress([]byte("test"))
	if err != nil {
		t.Errorf("Compress failed: %v", err)
		return
	}

	enc := encode.NewEncoding(encode.Base64)
	encodedStr, err := enc.Encode(data)
	if err != nil {
		t.Errorf("Encoding failed: %v", err)
		return
	}
	t.Log(encodedStr)

	decompressed, err := Decompress(data)
	if err != nil {
		t.Errorf("Decompress failed: %v", err)
		return
	}

	if !bytes.Equal(decompressed, []byte("test")) {
		t.Errorf("Decompressed data does not match original data")
		return
	}
}

func TestDecompress(t *testing.T) {
	enc := encode.NewEncoding(encode.Base64)
	data, err := enc.Decode([]byte("KLUv/QQAIQAAdGVzdDmBZ9s="))
	if err != nil {
		t.Error("unable to decode data")
		return
	}

	decompressed, err := Decompress(data)
	if err != nil {
		t.Errorf("Decompress failed: %v", err)
		return
	}

	if !bytes.Equal(decompressed, []byte("test")) {
		t.Errorf("Decompressed data does not match original data")
		return
	}
}
