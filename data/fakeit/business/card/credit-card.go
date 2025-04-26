package card

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateCreditCard() Card {
	now := time.Now()

	// Select random card brand
	brands := []string{"Visa", "Mastercard", "American Express", "Discover"}
	brand := brands[rand.Intn(len(brands))]

	// Get random spec for the selected brand
	specs := creditCardSpecs[brand]
	spec := specs[rand.Intn(len(specs))]

	// Generate card number
	number := generateNumber(spec.Prefix, spec.Length)

	// Generate expiry date (3-5 years from now)
	expiryYear := now.Year() + rand.Intn(3) + 3
	expiryMonth := rand.Intn(12) + 1

	// Generate CVV (4 digits for Amex, 3 for others)
	cvvLength := 3
	if brand == "American Express" {
		cvvLength = 4
	}
	cvv := generateCVV(cvvLength)

	return Card{
		Number:         formatNumber(number),
		CardholderName: "SAMPLE CARD",
		ExpiryMonth:    expiryMonth,
		ExpiryYear:     expiryYear,
		CVV:            cvv,
		Brand:          brand,
		IssueDate:      now,
	}
}

// Card types with their prefixes and lengths
var cardTypes = map[string]struct {
	prefixes []string
	length   int
}{
	"VISA":       {prefixes: []string{"4"}, length: 16},
	"MasterCard": {prefixes: []string{"51", "52", "53", "54", "55"}, length: 16},
	"Amex":       {prefixes: []string{"34", "37"}, length: 15},
}

func GenerateRandomCreditCard() Credicard {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Select random card type
	types := []string{"VISA", "MasterCard", "Amex"}
	cardType := types[rand.Intn(len(types))]

	// Generate card number
	number := generateValidCardNumber(cardType)

	// Generate expiry date (current year + 1-5 years)
	currentYear := time.Now().Year() % 100 // Get last two digits
	year := currentYear + rand.Intn(5) + 1
	month := rand.Intn(12) + 1
	expiry := fmt.Sprintf("%02d/%02d", month, year)

	// Generate CVV (3 digits for VISA/MC, 4 for Amex)
	cvvLength := 3
	if cardType == "Amex" {
		cvvLength = 4
	}
	cvv := ""
	for i := 0; i < cvvLength; i++ {
		cvv += strconv.Itoa(rand.Intn(10))
	}

	return Credicard{
		Number: number,
		Type:   cardType,
		Name:   "SAMPLE CARD", // Placeholder name
		Expiry: expiry,
		CVV:    cvv,
	}
}

func generateValidCardNumber(cardType string) string {
	cardInfo := cardTypes[cardType]
	prefix := cardInfo.prefixes[rand.Intn(len(cardInfo.prefixes))]
	length := cardInfo.length

	// Generate random digits for the remaining length
	remainingLength := length - len(prefix) - 1 // -1 for check digit
	number := prefix
	for i := 0; i < remainingLength; i++ {
		number += strconv.Itoa(rand.Intn(10))
	}

	// Calculate and append check digit
	checkDigit := generateCheckDigit(number)
	number += strconv.Itoa(checkDigit)

	return number
}

func generateCheckDigit(partial string) int {
	// Calculate sum according to Luhn algorithm
	sum := 0
	double := false

	for i := len(partial) - 1; i >= 0; i-- {
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

	// Calculate check digit
	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}
