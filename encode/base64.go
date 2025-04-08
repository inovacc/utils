package encode

import "encoding/base64"

type base64Encoding struct{}

func (b *base64Encoding) EncodeStr(s string) (string, error) {
	data, err := b.Encode([]byte(s))
	return string(data), err
}

func (b *base64Encoding) DecodeStr(s string) (string, error) {
	data, err := b.Decode([]byte(s))
	return string(data), err
}

func (b *base64Encoding) Encode(data []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(data)), nil
}

func (b *base64Encoding) Decode(data []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(data))
}
