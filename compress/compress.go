package compress

import (
	"github.com/inovacc/utils/v2/compress/brotli"
	"github.com/inovacc/utils/v2/compress/gzip"
	"github.com/inovacc/utils/v2/compress/lz4"
	"github.com/inovacc/utils/v2/compress/snappy"
	"github.com/inovacc/utils/v2/compress/zip"
	"github.com/inovacc/utils/v2/compress/zlib"
	"github.com/inovacc/utils/v2/compress/zstd"
)

type TypeStr string

const (
	TypeZstd   TypeStr = "zstd"
	TypeGzip   TypeStr = "gzip"
	TypeSnappy TypeStr = "snappy"
	TypeLz4    TypeStr = "lz4"
	TypeBrotli TypeStr = "brotli"
	TypeZlib   TypeStr = "zlib"
	TypeZip    TypeStr = "zip"
)

type Compress struct {
	Type TypeStr
}

func NewCompress(t TypeStr) *Compress {
	return &Compress{
		Type: t,
	}
}

func (c *Compress) Compress(data []byte) ([]byte, error) {
	switch c.Type {
	case TypeZstd:
		return zstd.Compress(data)
	case TypeGzip:
		return gzip.Compress(data)
	case TypeSnappy:
		return snappy.Compress(data)
	case TypeLz4:
		return lz4.Compress(data)
	case TypeBrotli:
		return brotli.Compress(data)
	case TypeZlib:
		return zlib.Compress(data)
	case TypeZip:
		return zip.Compress(data)
	default:
		return nil, nil
	}
}

func (c *Compress) Decompress(data []byte) ([]byte, error) {
	switch c.Type {
	case TypeZstd:
		return zstd.Decompress(data)
	case TypeGzip:
		return gzip.Decompress(data)
	case TypeSnappy:
		return snappy.Decompress(data)
	case TypeLz4:
		return lz4.Decompress(data)
	case TypeBrotli:
		return brotli.Decompress(data)
	case TypeZlib:
		return zlib.Decompress(data)
	case TypeZip:
		return zip.Decompress(data)
	default:
		return nil, nil
	}
}

func (c *Compress) String() string {
	switch c.Type {
	case TypeZstd:
		return "zstd"
	case TypeGzip:
		return "gzip"
	case TypeSnappy:
		return "snappy"
	case TypeLz4:
		return "lz4"
	case TypeBrotli:
		return "brotli"
	case TypeZlib:
		return "zlib"
	case TypeZip:
		return "zip"
	default:
		return ""
	}
}

func (c *Compress) TypeString() string {
	switch c.Type {
	case TypeZstd:
		return "zstd"
	case TypeGzip:
		return "gzip"
	case TypeSnappy:
		return "snappy"
	case TypeLz4:
		return "lz4"
	case TypeBrotli:
		return "brotli"
	case TypeZlib:
		return "zlib"
	case TypeZip:
		return "zip"
	default:
		return ""
	}
}
