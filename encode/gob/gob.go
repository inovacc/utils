package gob

import (
	"bytes"
	"encoding/gob"
)

func EncodeGob(data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}
	return nil, nil
}

func DecodeGob(data []byte, v any) error {
	buf := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buf).Decode(v); err != nil {
		return err
	}
	return nil
}
