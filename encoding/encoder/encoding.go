package encoder

// BaseType defines supported base encoding formats.
type BaseType int

const (
	Base58 BaseType = iota // Base58 encoding (commonly used in Bitcoin addresses)
	Base62                 // Base62 encoding (compact, URL-safe)
	Base64                 // Base64 encoding (standard in web/data transfer)
)

// Encoding defines a common interface for encoding and decoding operations.
// Implementations must support both byte-level and string-level operations.
type Encoding interface {
	Encode([]byte) ([]byte, error)    // Encode input bytes to encoded bytes
	EncodeStr(string) (string, error) // Encode input string to encoded string
	Decode([]byte) ([]byte, error)    // Decode encoded bytes to original bytes
	DecodeStr(string) (string, error) // Decode encoded string to original string
}

// NewEncoding returns an Encoding implementation based on the selected BaseType.
// It supports Base58, Base62, and Base64 encodings. Returns nil if base is unsupported.
func NewEncoding(base BaseType) Encoding {
	switch base {
	case Base58:
		return newBase58Encoding()
	case Base62:
		return newBase62Encoding()
	case Base64:
		return &base64Encoding{}
	default:
		return nil
	}
}
