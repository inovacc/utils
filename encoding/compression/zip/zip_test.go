package zip

import (
	"bytes"
	"testing"

	"github.com/inovacc/utils/v2/encoding/encoder"
)

func TestCompress(t *testing.T) {
	data, err := Compress([]byte("test"))
	if err != nil {
		t.Errorf("Compress failed: %v", err)
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

func TestDecompress(t *testing.T) {
	enc := encoder.NewEncoding(encoder.Base64)
	data, err := enc.Decode([]byte("UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAEAAAAZGF0YSpJLS4BBAAA//9QSwcIDH5/2AoAAAAEAAAAUEsBAhQAFAAIAAgAAAAAAAx+f9gKAAAABAAAAAQAAAAAAAAAAAAAAAAAAAAAAGRhdGFQSwUGAAAAAAEAAQAyAAAAPAAAAAAA"))
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
