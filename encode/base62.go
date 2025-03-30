package encode

import (
	"github.com/inovacc/base62"
)

func Base62Encode(data []byte) string {
	return base62.Encode(data)
}

func Base62Decode(data string) ([]byte, error) {
	decoded, err := base62.Decode(data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
