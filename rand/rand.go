package rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

func RandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, fmt.Errorf("invalid range: min (%d) must be less than max (%d)", min, max)
	}
	diff := max - min
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	return int(num.Int64()) + min, nil
}

func RandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return b, nil
}
