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
	Data []byte
}

func NewCompress(t TypeStr, data []byte) *Compress {
	return &Compress{
		Type: t,
		Data: data,
	}
}

func (c *Compress) Compress() ([]byte, error) {
	switch c.Type {
	case TypeZstd:
		return zstd.Compress(c.Data)
	case TypeGzip:
		return gzip.Compress(c.Data)
	case TypeSnappy:
		return snappy.Compress(c.Data)
	case TypeLz4:
		return lz4.Compress(c.Data)
	case TypeBrotli:
		return brotli.Compress(c.Data)
	case TypeZlib:
		return zlib.Compress(c.Data)
	case TypeZip:
		return zip.Compress(c.Data)
	default:
		return nil, nil
	}
}

func (c *Compress) Decompress() ([]byte, error) {
	switch c.Type {
	case TypeZstd:
		return zstd.Decompress(c.Data)
	case TypeGzip:
		return gzip.Decompress(c.Data)
	case TypeSnappy:
		return snappy.Decompress(c.Data)
	case TypeLz4:
		return lz4.Decompress(c.Data)
	case TypeBrotli:
		return brotli.Decompress(c.Data)
	case TypeZlib:
		return zlib.Decompress(c.Data)
	case TypeZip:
		return zip.Decompress(c.Data)
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
