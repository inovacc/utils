package card

import (
	"math/rand"
	"time"
)

func GenerateDebitCard() Card {
	now := time.Now()

	// Select random card brand
	brands := []string{"Visa Electron", "Maestro", "Visa Debit", "Mastercard Debit"}
	brand := brands[rand.Intn(len(brands))]

	// Get random spec for the selected brand
	specs := debitCardSpecs[brand]
	spec := specs[rand.Intn(len(specs))]

	// Generate card number
	number := generateNumber(spec.Prefix, spec.Length)

	// Generate expiry date (2-4 years from now, as debit cards typically have shorter validity)
	expiryYear := now.Year() + rand.Intn(3) + 2
	expiryMonth := rand.Intn(12) + 1

	// Generate CVV (always 3 digits for debit cards)
	cvv := generateCVV(3)

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
