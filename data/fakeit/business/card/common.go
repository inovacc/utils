package card

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	Number         string
	CardholderName string
	ExpiryMonth    int
	ExpiryYear     int
	CVV            string
	Brand          string
	IssueDate      time.Time
}

type CardDetails struct {
	Prefix string
	Length int
}

// CardSpec defines the specification for a card type
type CardSpec struct {
	prefixes []string
	lengths  []int
}

// Card specifications based on provided information
var cardSpecs = map[string]CardSpec{
	"Visa": {
		prefixes: []string{"4"},
		lengths:  []int{13, 16, 19},
	},
	"Mastercard": {
		prefixes: generatePrefixRange(51, 55, append(generatePrefixRange(2221, 2720))),
		lengths:  []int{16},
	},
	"American Express": {
		prefixes: []string{"34", "37"},
		lengths:  []int{15},
	},
	"Discover": {
		prefixes: append(
			append(
				[]string{"6011"},
				generatePrefixRange(622126, 622925)...,
			),
			append(
				generatePrefixRange(644, 649),
				"65",
			)...,
		),
		lengths: []int{16},
	},
	"Diners Club": {
		prefixes: append(
			generatePrefixRange(300, 305),
			[]string{"36", "38"}...,
		),
		lengths: []int{14},
	},
	"JCB": {
		prefixes: generatePrefixRange(3528, 3589),
		lengths:  []int{16, 17, 18, 19},
	},
}

var creditCardSpecs = map[string][]CardDetails{
	"Visa": {{Prefix: "4", Length: 16}},
	"Mastercard": {
		{Prefix: "51", Length: 16},
		{Prefix: "52", Length: 16},
		{Prefix: "53", Length: 16},
		{Prefix: "54", Length: 16},
		{Prefix: "55", Length: 16},
	},
	"American Express": {
		{Prefix: "34", Length: 15},
		{Prefix: "37", Length: 15},
	},
	"Discover": {
		{Prefix: "6011", Length: 16},
		{Prefix: "65", Length: 16},
	},
}

var debitCardSpecs = map[string][]CardDetails{
	"Visa Electron": {
		{Prefix: "4026", Length: 16},
		{Prefix: "417500", Length: 16},
		{Prefix: "4508", Length: 16},
		{Prefix: "4844", Length: 16},
		{Prefix: "4913", Length: 16},
		{Prefix: "4917", Length: 16},
	},
	"Maestro": {
		{Prefix: "5018", Length: 16},
		{Prefix: "5020", Length: 16},
		{Prefix: "5038", Length: 16},
		{Prefix: "6304", Length: 16},
		{Prefix: "6759", Length: 16},
		{Prefix: "6761", Length: 16},
		{Prefix: "6763", Length: 16},
	},
	"Visa Debit": {{Prefix: "4", Length: 16}},
	"Mastercard Debit": {
		{Prefix: "51", Length: 16},
		{Prefix: "52", Length: 16},
		{Prefix: "53", Length: 16},
		{Prefix: "54", Length: 16},
		{Prefix: "55", Length: 16},
	},
}

// Helper function to generate prefix ranges
func generatePrefixRange(start, end int, additional ...string) []string {
	var prefixes []string
	for i := start; i <= end; i++ {
		prefixes = append(prefixes, fmt.Sprintf("%d", i))
	}
	return append(prefixes, additional...)
}

func GenerateCard() BankCard {
	rand.Seed(time.Now().UnixNano())

	// Get random card type
	types := make([]string, 0, len(cardSpecs))
	for cardType := range cardSpecs {
		types = append(types, cardType)
	}
	cardType := types[rand.Intn(len(types))]

	// Generate card details
	number := generateValidCardNumber(cardType)

	// Generate expiry date (current year + 2-5 years)
	currentYear := time.Now().Year() % 100
	year := currentYear + rand.Intn(4) + 2
	month := rand.Intn(12) + 1
	expiry := fmt.Sprintf("%02d/%02d", month, year)

	// Generate CVV (4 digits for Amex, 3 for others)
	cvvLength := 3
	if cardType == "American Express" {
		cvvLength = 4
	}
	cvv := generateRandomDigits(cvvLength)

	return BankCard{
		Number: number,
		Type:   cardType,
		Name:   "SAMPLE CARD",
		Expiry: expiry,
		CVV:    cvv,
	}
}

func generateValidCardNumber(cardType string) string {
	spec := cardSpecs[cardType]

	// Select random prefix and length
	prefix := spec.prefixes[rand.Intn(len(spec.prefixes))]
	length := spec.lengths[rand.Intn(len(spec.lengths))]

	// Generate random digits for remaining length
	remainingLength := length - len(prefix) - 1 // -1 for check digit
	number := prefix + generateRandomDigits(remainingLength)

	// Calculate and append check digit
	checkDigit := generateCheckDigit(number)
	finalNumber := number + strconv.Itoa(checkDigit)

	return formatCardNumber(finalNumber, cardType)
}

func generateRandomDigits(length int) string {
	digits := make([]byte, length)
	for i := 0; i < length; i++ {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	return string(digits)
}

func formatCardNumber(number string, cardType string) string {
	var parts []string
	switch cardType {
	case "American Express":
		// Format: XXXX XXXXXX XXXXX
		parts = []string{
			number[:4],
			number[4:10],
			number[10:],
		}
	case "Diners Club":
		// Format: XXXX XXXXXX XX
		parts = []string{
			number[:4],
			number[4:10],
			number[10:],
		}
	default:
		// Standard format: XXXX XXXX XXXX XXXX (or longer for 19-digit cards)
		for i := 0; i < len(number); i += 4 {
			end := i + 4
			if end > len(number) {
				end = len(number)
			}
			parts = append(parts, number[i:end])
		}
	}
	return strings.Join(parts, " ")
}

func generateCheckDigit(partial string) int {
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

func cardBrand(number string) string {
	number = strings.ReplaceAll(number, " ", "") // Remove spaces

	switch {
	case strings.HasPrefix(number, "4"):
		return "Visa"
	case strings.HasPrefix(number, "34"), strings.HasPrefix(number, "37"):
		return "American Express"
	case strings.HasPrefix(number, "5"):
		return "Mastercard"
	case strings.HasPrefix(number, "6011"), strings.HasPrefix(number, "65"):
		return "Discover"
	case strings.HasPrefix(number, "35"):
		return "JCB"
	case strings.HasPrefix(number, "30"), strings.HasPrefix(number, "36"), strings.HasPrefix(number, "38"):
		return "Diners Club"
	default:
		return "Unknown Brand"
	}
}

func init() {
	// Card types with their prefixes and lengths
	var cardTypes = map[string]struct {
		prefixes []string
		length   int
	}{
		"VISA":       {prefixes: []string{"4"}, length: 16},
		"MasterCard": {prefixes: []string{"51", "52", "53", "54", "55"}, length: 16},
		"Amex":       {prefixes: []string{"34", "37"}, length: 15},
	}
}
