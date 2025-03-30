package argon2

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/inovacc/base58"
	"github.com/inovacc/utils/v2/rand"
	"golang.org/x/crypto/argon2"
)

type Params struct {
	Memory      uint32 `json:"m"`
	Iterations  uint32 `json:"i"`
	Parallelism uint8  `json:"p"`
	SaltLength  uint32 `json:"s"`
	KeyLength   uint32 `json:"k"`
}

type Hash struct {
	Data   []byte  `json:"data"`
	Salt   []byte  `json:"salt"`
	Params *Params `json:"params"`
}

var params *Params

func init() {
	params = &Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func hashPassword(password string, salt []byte, p *Params) ([]byte, error) {
	hash := argon2.IDKey([]byte(password),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)
	return hash, nil
}

// HashPassword generates a secure Argon2ID hash for the given password and parameters.
func HashPassword(password string, p *Params) (string, error) {
	if p == nil {
		p = params
	}

	if p.Memory < 1 || p.Iterations < 1 || p.Parallelism < 1 || p.SaltLength < 1 || p.KeyLength < 1 {
		return "", errors.New("invalid parameters")
	}

	if len(password) == 0 {
		return "", errors.New("password cannot be empty")
	}

	salt := rand.RandomBytes(p.SaltLength)
	hash, err := hashPassword(password, salt, p)
	if err != nil {
		return "", err
	}
	stored := &Hash{
		Data:   hash,
		Salt:   salt,
		Params: p,
	}
	encoded, err := json.Marshal(stored)
	if err != nil {
		return "", err
	}
	return base58.StdEncoding.Encode(encoded), nil
}

// CheckPasswordHash compares a plain-text password with a stored hash (JSON encoded).
func CheckPasswordHash(encoded, password string) (bool, error) {
	var stored Hash
	decode, err := base58.StdEncoding.Decode(encoded)
	if err != nil {
		return false, errors.New("invalid encoding hash")
	}

	if len(decode) == 0 {
		return false, errors.New("invalid hash length")
	}

	if err := json.Unmarshal(decode, &stored); err != nil {
		return false, errors.New("invalid unmarshal hash")
	}

	hash, err := hashPassword(password, stored.Salt, stored.Params)
	if err != nil {
		return false, err
	}
	return bytes.Equal(hash, stored.Data), nil
}
