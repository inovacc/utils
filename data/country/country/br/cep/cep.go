package cep

import (
	"fmt"
	"regexp"
	"strings"
)

// IsValidCEP validates a Brazilian postal code (CEP).
func IsValidCEP(input string) bool {
	input = UnformatCEP(input)

	match, _ := regexp.MatchString(`^\d{8}$`, input)
	return match
}

func FormatCEP(input string) string {
	input = UnformatCEP(input)

	if len(input) != 8 {
		return ""
	}

	return fmt.Sprintf("%s-%s", input[:5], input[5:])
}

func UnformatCEP(cep string) string {
	clean := strings.ReplaceAll(cep, "-", "")
	clean = strings.ReplaceAll(clean, ".", "")
	clean = strings.TrimSpace(clean)

	if len(clean) != 8 {
		return ""
	}

	return clean
}
