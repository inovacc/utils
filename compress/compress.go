package compress

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
		return ZstdCompress(c.Data)
	case TypeGzip:
		return GzipCompress(c.Data)
	case TypeSnappy:
		return SnappyCompress(c.Data)
	case TypeLz4:
		return Lz4Compress(c.Data)
	case TypeBrotli:
		return BrotliCompress(c.Data)
	case TypeZlib:
		return ZlibCompress(c.Data)
	case TypeZip:
		return ZipCompress(c.Data)
	default:
		return nil, nil
	}
}

func (c *Compress) Decompress() ([]byte, error) {
	switch c.Type {
	case TypeZstd:
		return ZstdDecompress(c.Data)
	case TypeGzip:
		return GzipDecompress(c.Data)
	case TypeSnappy:
		return SnappyDecompress(c.Data)
	case TypeLz4:
		return Lz4Decompress(c.Data)
	case TypeBrotli:
		return BrotliDecompress(c.Data)
	case TypeZlib:
		return ZlibDecompress(c.Data)
	case TypeZip:
		return ZipDecompress(c.Data)
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
