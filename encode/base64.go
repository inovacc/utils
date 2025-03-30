package encode

import (
	"encoding/base64"
	"github.com/inovacc/base62"
)

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	decoded, err := base62.Decode(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
