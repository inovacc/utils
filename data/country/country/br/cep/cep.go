package cep

import (
	"regexp"
	"strings"
)

// IsValidCEP validates a Brazilian postal code (CEP).
func IsValidCEP(input string) bool {
	clean := strings.ReplaceAll(input, "-", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.TrimSpace(clean)

	match, _ := regexp.MatchString(`^\d{8}$`, clean)
	return match
}
