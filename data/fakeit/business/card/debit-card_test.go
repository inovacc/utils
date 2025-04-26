package card

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Card types with their prefixes and lengths
var cardTypes = map[string]struct {
	prefixes []string
	length   int
}{
	"Visa Debit":       {prefixes: []string{"4"}, length: 16},
	"Mastercard Debit": {prefixes: []string{"51", "52", "53", "54", "55", "2221", "2720"}, length: 16},
	"Maestro":          {prefixes: []string{"5018", "5020", "5038", "6304", "6759", "6761", "6763"}, length: 16},
	"Visa Electron":    {prefixes: []string{"4026", "417500", "4508", "4844", "4913", "4917"}, length: 16},
}

// Rename the struct to better represent both credit and debit cards
type DebitCard struct {
	Number      string
	Type        string
	Name        string
	Expiry      string
	CVV         string
	BankName    string
	IsDebitCard bool
}

func GenerateRandomDebitCard() DebitCard {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Select random debit card type
	types := []string{"Visa Debit", "Mastercard Debit", "Maestro", "Visa Electron"}
	cardType := types[rand.Intn(len(types))]

	// List of sample bank names
	bankNames := []string{
		"National Bank",
		"City Bank",
		"Community Bank",
		"Metro Bank",
		"Regional Bank",
	}

	// Generate card number
	number := generateValidCardNumber(cardType)

	// Generate expiry date (current year + 2-4 years, as debit cards typically have shorter validity)
	currentYear := time.Now().Year() % 100 // Get last two digits
	year := currentYear + rand.Intn(3) + 2
	month := rand.Intn(12) + 1
	expiry := fmt.Sprintf("%02d/%02d", month, year)

	// Generate CVV (always 3 digits for debit cards)
	cvv := ""
	for i := 0; i < 3; i++ {
		cvv += strconv.Itoa(rand.Intn(10))
	}

	return DebitCard{
		Number:      number,
		Type:        cardType,
		Name:        "SAMPLE CARD",
		Expiry:      expiry,
		CVV:         cvv,
		BankName:    bankNames[rand.Intn(len(bankNames))],
		IsDebitCard: true,
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

	return formatCardNumber(number, cardType)
}

func formatCardNumber(number string, cardType string) string {
	// Format the card number with spaces or hyphens based on card type
	var formatted string
	switch cardType {
	case "Maestro":
		// Maestro cards can have varying formats, using 4-4-4-4
		for i := 0; i < len(number); i += 4 {
			end := i + 4
			if end > len(number) {
				end = len(number)
			}
			if i > 0 {
				formatted += " "
			}
			formatted += number[i:end]
		}
	default:
		// Standard format: 4-4-4-4
		for i := 0; i < len(number); i += 4 {
			end := i + 4
			if end > len(number) {
				end = len(number)
			}
			if i > 0 {
				formatted += " "
			}
			formatted += number[i:end]
		}
	}
	return formatted
}

func generateCheckDigit(partial string) int {
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

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit
}
