package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// Option is a functional option type for configuring password generation.
type Option func(*Password)

// passwordOptions defines character sets for password generation.
var passwordOptions = map[string]string{
	"num":         "1234567890",
	"specialChar": "!@#$%&'()*+,^-./:;<=>?[]_`{~}|",
	"lowerCase":   "abcdefghijklmnopqrstuvwxyz",
	"upperCase":   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
}

// Password holds configuration for generating passwords.
type Password struct {
	length      int
	hasNumbers  bool
	hashSpecial bool
	hasLower    bool
	hasUpper    bool
}

// NewPassword constructs a Password with provided options.
// By default, it creates a password of length 8 with no character types enabled.
func NewPassword(opts ...Option) *Password {
	p := &Password{
		length: 8,
	}

	for _, opt := range opts {
		opt(p)
	}
	return p
}

// WithLength sets the length of the password (8–128).
func WithLength(length int) Option {
	return func(p *Password) {
		if length >= 8 && length <= 128 {
			p.length = length
		} else {
			fmt.Println("Length must be between 8 and 128 characters.")
		}
	}
}

// WithNumbers enables or disables numeric characters in the password.
func WithNumbers(enabled bool) Option {
	return func(p *Password) {
		p.hasNumbers = enabled
	}
}

// WithSpecial enables or disables special characters in the password.
func WithSpecial(enabled bool) Option {
	return func(p *Password) {
		p.hashSpecial = enabled
	}
}

// WithLower enables or disables lowercase letters in the password.
func WithLower(enabled bool) Option {
	return func(p *Password) {
		p.hasLower = enabled
	}
}

// WithUpper enables or disables uppercase letters in the password.
func WithUpper(enabled bool) Option {
	return func(p *Password) {
		p.hasUpper = enabled
	}
}

// Generate builds a password using the configured options.
// It ensures at least one character from each enabled category is included.
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

	// Fill the rest with the password
	for len(passChars) < p.length {
		passChars = append(passChars, p.getRandomChar(passInfo))
	}

	p.shuffleRunes(passChars)
	return string(passChars), nil
}

// getRandomChar selects one rune at random from the input string.
func (p *Password) getRandomChar(fromString string) rune {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(fromString))))
	return rune(fromString[index.Int64()])
}

// shuffleRunes performs Fisher–Yates shuffle to randomize rune order.
func (p *Password) shuffleRunes(runes []rune) {
	for i := len(runes) - 1; i > 0; i-- {
		jBig, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jBig.Int64())
		runes[i], runes[j] = runes[j], runes[i]
	}
}
