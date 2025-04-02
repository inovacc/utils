// Fonte : https://www.serpro.gov.br/menu/noticias/videos/calculodvcnpjalfanaumerico.pdf
//

package cnpj

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Character-to-numeric value mapping for alphanumeric CNPJ
var charToValue = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 17, 'B': 18, 'C': 19, 'D': 20, 'E': 21, 'F': 22, 'G': 23, 'H': 24, 'I': 25, 'J': 26,
	'K': 27, 'L': 28, 'M': 29, 'N': 30, 'O': 31, 'P': 32, 'Q': 33, 'R': 34, 'S': 35, 'T': 36,
	'U': 37, 'V': 38, 'W': 39, 'X': 40, 'Y': 41, 'Z': 42,
}

// calculateCheckDigit calculates a check digit using modulo 11
func calculateCheckDigit(cnpj string) int {
	weights := []int{2, 3, 4, 5, 6, 7, 8, 9}
	sum := 0
	j := 0

	// Iterate through CNPJ from right to left applying the weights
	for i := len(cnpj) - 1; i >= 0; i-- {
		val, ok := charToValue[rune(cnpj[i])]
		if !ok {
			return -1 // Invalid character
		}
		sum += val * weights[j]
		j = (j + 1) % len(weights) // Loop weights after the 8th element
	}

	remainder := sum % 11
	if remainder == 0 || remainder == 1 {
		return 0
	}
	return 11 - remainder
}

// GenerateCNPJ generates a valid alphanumeric CNPJ
func GenerateCNPJ() string {
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder

	// Generate the first 12 random characters (numbers or letters)
	for i := 0; i < 12; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte(byte('0' + rand.Intn(10))) // Number
		} else {
			sb.WriteByte(byte('A' + rand.Intn(26))) // Letter
		}
	}

	cnpjBase := sb.String()

	// Calculate the two check digits
	dv1 := calculateCheckDigit(cnpjBase)
	dv2 := calculateCheckDigit(cnpjBase + strconv.Itoa(dv1))

	return fmt.Sprintf("%s%d%d", cnpjBase, dv1, dv2)
}

// ValidateCNPJ checks whether an alphanumeric CNPJ is valid
func ValidateCNPJ(cnpj string) bool {
	// Remove formatting
	cnpj = UnformatCNPJ(cnpj)

	if len(cnpj) != 14 {
		return false
	}

	base := cnpj[:12]
	dv1, _ := strconv.Atoi(string(cnpj[12]))
	dv2, _ := strconv.Atoi(string(cnpj[13]))

	return calculateCheckDigit(base) == dv1 && calculateCheckDigit(base+strconv.Itoa(dv1)) == dv2
}

// FormatCNPJ formats an alphanumeric CNPJ in the pattern "12.ABC.345/01DE-XX"
func FormatCNPJ(cnpj string) string {
	cnpj = UnformatCNPJ(cnpj)

	if len(cnpj) != 14 {
		return "Invalid CNPJ"
	}

	// CNPJ mask: "XX.XXX.XXX/XXXX-XX"
	mask := "##.###.###/####-##"
	result := make([]rune, len(mask))
	cnpjIndex := 0

	for i, r := range mask {
		if r == '#' {
			result[i] = rune(cnpj[cnpjIndex])
			cnpjIndex++
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// UnformatCNPJ removes formatting from an alphanumeric CNPJ
func UnformatCNPJ(cnpj string) string {
	re := regexp.MustCompile(`[^0-9A-Z]`)
	return strings.ToUpper(re.ReplaceAllString(cnpj, ""))
}
