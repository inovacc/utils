package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type Option func(*Password)

var passwordOptions = map[string]string{
	"num":         "1234567890",
	"specialChar": "!@#$%&'()*+,^-./:;<=>?[]_`{~}|",
	"lowerCase":   "abcdefghijklmnopqrstuvwxyz",
	"upperCase":   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}

type Password struct {
	length      int
	hasNumbers  bool
	hashSpecial bool
	hasLower    bool
	hasUpper    bool
}

func NewPassword(opts ...Option) *Password {
	p := &Password{
		length: 8,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithLength(length int) Option {
	return func(p *Password) {
		if length >= 8 && length <= 128 {
			p.length = length
		} else {
			fmt.Println("Length must be between 8 and 128 characters.")
		}
	}
}

func WithNumbers(enabled bool) Option {
	return func(p *Password) {
		p.hasNumbers = enabled
	}
}

func WithSpecial(enabled bool) Option {
	return func(p *Password) {
		p.hashSpecial = enabled
	}
}

func WithLower(enabled bool) Option {
	return func(p *Password) {
		p.hasLower = enabled
	}
}

func WithUpper(enabled bool) Option {
	return func(p *Password) {
		p.hasUpper = enabled
	}
}

func (p *Password) Generate() (string, error) {
	var passInfo string
	var passChars []rune

	if p.length < 8 || p.length > 128 {
		return "", errors.New("length must be between 8 and 128 characters")
	}

	if p.hasNumbers {
		passInfo += passwordOptions["num"]
		passChars = append(passChars, p.getRandomChar(passwordOptions["num"]))
	}

	if p.hashSpecial {
		passInfo += passwordOptions["specialChar"]
		passChars = append(passChars, p.getRandomChar(passwordOptions["specialChar"]))
	}

	if p.hasLower {
		passInfo += passwordOptions["lowerCase"]
		passChars = append(passChars, p.getRandomChar(passwordOptions["lowerCase"]))
	}

	if p.hasUpper {
		passInfo += passwordOptions["upperCase"]
		passChars = append(passChars, p.getRandomChar(passwordOptions["upperCase"]))
	}

	if passInfo == "" {
		return "", errors.New("no character sets selected, please enable at least one")
	}

	for len(passChars) < p.length {
		passChars = append(passChars, p.getRandomChar(passInfo))
	}

	p.shuffleRunes(passChars)
	return string(passChars), nil
}

func (p *Password) getRandomChar(fromString string) rune {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(fromString))))
	return rune(fromString[index.Int64()])
}

func (p *Password) shuffleRunes(runes []rune) {
	for i := len(runes) - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
}
