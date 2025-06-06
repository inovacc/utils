package card

import (
	"fmt"
	"testing"
)

func TestGenerateCreditCard(t *testing.T) {
	// Generate a credit card
	creditCard := GenerateCreditCard()
	fmt.Printf("Credit Card:\n")
	fmt.Printf("Brand: %s\n", creditCard.Brand)
	fmt.Printf("Number: %s\n", creditCard.Number)
	fmt.Printf("Name: %s\n", creditCard.CardholderName)
	fmt.Printf("Expiry: %02d/%d\n", creditCard.ExpiryMonth, creditCard.ExpiryYear)
	fmt.Printf("CVV: %s\n", creditCard.CVV)
	fmt.Printf("Issue Date: %v\n\n", creditCard.IssueDate.Format("01/2006"))
}
