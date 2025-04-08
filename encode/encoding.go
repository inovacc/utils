package encode

type BaseType int

const (
	Base58 BaseType = iota
	Base62
	Base64
)

type Encoding interface {
	Encode([]byte) ([]byte, error)
	EncodeStr(string) (string, error)

	Decode([]byte) ([]byte, error)
	DecodeStr(string) (string, error)
}

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
