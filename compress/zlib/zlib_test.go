package zlib

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

	encodedStr := encode.Base64Encode(data)
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
	data, err := encode.Base64Decode("eJwqSS0uAQQAAP//BF0BwQ==")
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
