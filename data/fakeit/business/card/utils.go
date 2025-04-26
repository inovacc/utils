package card

import (
	"math/rand"
	"strconv"
	"strings"
)

func generateNumber(prefix string, length int) string {
	// Generate random digits for remaining length
	remainingLength := length - len(prefix) - 1 // -1 for check digit
	number := prefix
	for i := 0; i < remainingLength; i++ {
		number += strconv.Itoa(rand.Intn(10))
	}

	// Calculate and append check digit
	checkDigit := generateLuhnCheckDigit(number)
	return number + strconv.Itoa(checkDigit)
}

func generateLuhnCheckDigit(partial string) int {
	sum := 0
	double := len(partial)%2 == 0

	for i := 0; i < len(partial); i++ {
		digit := int(partial[i] - '0')
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}

	return (10 - (sum % 10)) % 10
}

func generateCVV(length int) string {
	cvv := ""
	for i := 0; i < length; i++ {
		cvv += strconv.Itoa(rand.Intn(10))
	}
	return cvv
}

func formatNumber(number string) string {
	var parts []string
	for i := 0; i < len(number); i += 4 {
		end := i + 4
		if end > len(number) {
			end = len(number)
		}
		parts = append(parts, number[i:end])
	}
	return strings.Join(parts, " ")
}
