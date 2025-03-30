package encode

import (
	"github.com/inovacc/base58"
)

func Base58Encode(data []byte) string {
	return base58.StdEncoding.EncodeToString(data)
}

func Base58Decode(data string) ([]byte, error) {
	decoded, err := base58.StdEncoding.Decode(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
