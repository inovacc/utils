package rand

import (
	"crypto/rand"
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

func RandomInt(min, max int) int {
	if min >= max {
		panic("invalid range")
	}
	diff := max - min
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	return int(num.Int64()) + min
}

func RandomBytes(n uint32) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
