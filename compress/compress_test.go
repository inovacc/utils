package compress

import (
	"bytes"
	"testing"
)

func TestNewCompress(t *testing.T) {
	v := NewCompress(TypeZip, []byte("test"))
	b1, err := v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 := NewCompress(TypeZip, b1)
	b2, err := v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeZip {
		t.Errorf("Expected TypeZip, got %s", v.Type)
		return
	}

	if v2.Type != TypeZip {
		t.Errorf("Expected TypeZip, got %s", v2.Type)
		return
	}

	v = NewCompress(TypeGzip, []byte("test"))
	b1, err = v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 = NewCompress(TypeGzip, b1)
	b2, err = v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeGzip {
		t.Errorf("Expected TypeGzip, got %s", v.Type)
		return
	}

	if v2.Type != TypeGzip {
		t.Errorf("Expected TypeGzip, got %s", v2.Type)
		return
	}

	v = NewCompress(TypeSnappy, []byte("test"))
	b1, err = v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 = NewCompress(TypeSnappy, b1)
	b2, err = v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeSnappy {
		t.Errorf("Expected TypeSnappy, got %s", v.Type)
		return
	}

	if v2.Type != TypeSnappy {
		t.Errorf("Expected TypeSnappy, got %s", v2.Type)
		return
	}

	v = NewCompress(TypeLz4, []byte("test"))
	b1, err = v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 = NewCompress(TypeLz4, b1)
	b2, err = v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeLz4 {
		t.Errorf("Expected TypeLz4, got %s", v.Type)
		return
	}

	if v2.Type != TypeLz4 {
		t.Errorf("Expected TypeLz4, got %s", v2.Type)
		return
	}

	v = NewCompress(TypeBrotli, []byte("test"))
	b1, err = v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 = NewCompress(TypeBrotli, b1)
	b2, err = v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeBrotli {
		t.Errorf("Expected TypeBrotli, got %s", v.Type)
		return
	}

	if v2.Type != TypeBrotli {
		t.Errorf("Expected TypeBrotli, got %s", v2.Type)
		return
	}

	v = NewCompress(TypeZlib, []byte("test"))
	b1, err = v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 = NewCompress(TypeZlib, b1)
	b2, err = v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeZlib {
		t.Errorf("Expected TypeZlib, got %s", v.Type)
		return
	}

	if v2.Type != TypeZlib {
		t.Errorf("Expected TypeZlib, got %s", v2.Type)
		return
	}

	v = NewCompress(TypeZstd, []byte("test"))
	b1, err = v.Compress()
	if err != nil {
		t.Errorf("Error compressing: %v", err)
		return
	}

	v2 = NewCompress(TypeZstd, b1)
	b2, err = v2.Decompress()
	if err != nil {
		t.Errorf("Error decompressing: %v", err)
		return
	}

	if !bytes.Equal(b2, []byte("test")) {
		t.Errorf("Expected 'test', got '%s'", string(b2))
		return
	}

	if v.Type != TypeZstd {
		t.Errorf("Expected TypeZstd, got %s", v.Type)
		return
	}

	if v2.Type != TypeZstd {
		t.Errorf("Expected TypeZstd, got %s", v2.Type)
		return
	}
}
