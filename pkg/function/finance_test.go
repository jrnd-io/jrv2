// Copyright © 2024 JR team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package function_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	// Define test cases
	testCases := []struct {
		length int
	}{
		{5},
		{10},
		{15},
		{20},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Length: %d", tc.length), func(t *testing.T) {
			result := function.Account(tc.length)

			// Check if the length of the result is correct
			if len(result) != tc.length {
				t.Errorf("Expected length %d, but got %d", tc.length, len(result))
			}

			// Check if the result contains only digits
			for _, char := range result {
				if !unicode.IsDigit(char) {
					t.Errorf("Expected only digits, but got non-digit character: %c", char)
				}
			}
		})
	}

}

func TestAmount(t *testing.T) {
	// Define test cases
	testCases := []struct {
		min      float32
		max      float32
		currency string
	}{
		{0.0, 10.0, "$"},
		{5.0, 15.0, "€"},
		{10.0, 20.0, "£"},
		{100.0, 200.0, "¥"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Min: %f, Max: %f, Currency: %s", tc.min, tc.max, tc.currency), func(t *testing.T) {
			result := function.Amount(tc.min, tc.max, tc.currency)

			// Check if the result contains the currency symbol
			if !strings.Contains(result, tc.currency) {
				t.Errorf("Expected currency %s, but got %s", tc.currency, result)
			}

			// Extract the amount part from the result
			amountStr := strings.TrimPrefix(result, tc.currency)
			amount, err := strconv.ParseFloat(amountStr, 32)
			if err != nil {
				t.Errorf("Failed to parse amount: %v", err)
			}

			// Check if the amount is within the range
			if float32(amount) < tc.min || float32(amount) > tc.max {
				t.Errorf("Expected amount between %f and %f, but got %f", tc.min, tc.max, amount)
			}
		})
	}
}

func TestBitcoin(t *testing.T) {
	// Define the expected Bitcoin address pattern
	pattern := `^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`
	re := regexp.MustCompile(pattern)

	// Generate a Bitcoin address
	result := function.Bitcoin()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated Bitcoin address %s does not match the expected pattern", result)
	}
}
func TestCusip(t *testing.T) {
	// Define the expected CUSIP pattern
	pattern := `^[0-9]{3}[0-9A-Z]{5}[0-9]$`
	re := regexp.MustCompile(pattern)

	// Generate a CUSIP number
	result := function.Cusip()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated CUSIP %s does not match the expected pattern", result)
	}

	// Verify the check digit
	cusipWithoutCheckDigit := result[:8]
	expectedCheckDigit := function.CusipCheckDigit(cusipWithoutCheckDigit)
	actualCheckDigit := result[8:]

	if expectedCheckDigit != actualCheckDigit {
		t.Errorf("Expected check digit %s, but got %s", expectedCheckDigit, actualCheckDigit)
	}
}

func TestCreditCardCVV(t *testing.T) {
	// Define test cases
	testCases := []struct {
		length int
	}{
		{3},
		{4},
		{5},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Length: %d", tc.length), func(t *testing.T) {
			result := function.CreditCardCVV(tc.length)

			// Check if the result has the correct length
			if len(result) != tc.length {
				t.Errorf("Expected length %d, but got %d", tc.length, len(result))
			}

			// Check if the result contains only digits
			for _, char := range result {
				if !unicode.IsDigit(char) {
					t.Errorf("Expected only digits, but got %c", char)
				}
			}
		})
	}
}

func TestCreditCard(t *testing.T) {
	// Define test cases
	testCases := []struct {
		issuer  string
		pattern string
	}{
		{"visa", "^4\\d{12}(?:\\d{2})?$"},
		{"mastercard", "^5[1-5]\\d{13}$"},
		{"amex", "^3[47]\\d{12,13}$"},
		{"discover", "^(?:6011\\d{12})|(?:65\\d{13})$"},
	}

	for _, tc := range testCases {
		t.Run(tc.issuer, func(t *testing.T) {
			result := function.CreditCard(tc.issuer)

			// Check if the result matches the pattern
			re := regexp.MustCompile(tc.pattern)
			if !re.MatchString(result[:len(result)-1]) {
				t.Errorf("Generated credit card number %s does not match the expected pattern for %s", result, tc.issuer)
			}

			// Verify the Luhn check digit
			cardWithoutCheckDigit := result[:len(result)-1]
			expectedCheckDigit := function.LuhnCheckDigit(cardWithoutCheckDigit)
			actualCheckDigit := result[len(result)-1:]

			if expectedCheckDigit != actualCheckDigit {
				t.Errorf("Expected check digit %s, but got %s", expectedCheckDigit, actualCheckDigit)
			}
		})
	}
}

func TestEthereum(t *testing.T) {
	// Define the expected Ethereum address pattern
	pattern := `^0x[a-fA-F0-9]{40}$`
	re := regexp.MustCompile(pattern)

	// Generate an Ethereum address
	result := function.Ethereum()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated Ethereum address %s does not match the expected pattern", result)
	}
}

func TestIsin(t *testing.T) {
	// Define test cases
	testCases := []struct {
		country string
	}{
		{"US"},
		{"GB"},
		{"DE"},
	}

	for _, tc := range testCases {
		t.Run(tc.country, func(t *testing.T) {
			result := function.Isin(tc.country)

			// Define the expected ISIN pattern
			pattern := `^[A-Z]{2}[0-9A-Z]{9}[0-9]$`
			re := regexp.MustCompile(pattern)

			// Check if the result matches the pattern
			if !re.MatchString(result) {
				t.Errorf("Generated ISIN %s does not match the expected pattern", result)
			}

			// Verify the check digit
			isinWithoutCheckDigit := result[:11]
			expectedCheckDigit := function.IsinCheckDigit(isinWithoutCheckDigit)
			actualCheckDigit := result[11:]

			if expectedCheckDigit != actualCheckDigit {
				t.Errorf("Expected check digit %s, but got %s", expectedCheckDigit, actualCheckDigit)
			}
		})
	}
}

func TestSedol(t *testing.T) {
	// Define the expected SEDOL pattern
	pattern := `^[0-9BCDFGHJKLMNPQRSTVWXYZ]{6}[0-9]$`
	re := regexp.MustCompile(pattern)

	// Generate a SEDOL code
	result := function.Sedol()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated SEDOL %s does not match the expected pattern", result)
	}

	// Verify the check digit
	sedolWithoutCheckDigit := result[:6]
	expectedCheckDigit := function.SedolCheckDigit(sedolWithoutCheckDigit)
	actualCheckDigit := result[6:]

	if expectedCheckDigit != actualCheckDigit {
		t.Errorf("Expected check digit %s, but got %s", expectedCheckDigit, actualCheckDigit)
	}
}

func TestStockSymbol(t *testing.T) {
	emitter.GetState().Locale = TestLocale
	function.ClearCache(function.StockSymbolMap)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", function.StockSymbolMap), function.StockSymbolMap)
	if err != nil {
		t.Error(err)
		return
	}
	c := function.StockSymbol()
	assert.Contains(t, function.GetCache(function.StockSymbolMap), c)
	function.ClearCache(function.StockSymbolMap)
}

func TestSwift(t *testing.T) {
	// Define the expected SWIFT code pattern
	pattern := `^[A-Z]{4}[A-Z]{2}\d{2}\d{3}$`
	re := regexp.MustCompile(pattern)

	// Generate a SWIFT code
	result := function.Swift()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated SWIFT code %s does not match the expected pattern", result)
	}
}

func TestValor(t *testing.T) {
	// Define the expected pattern for Valor
	pattern := `^[0-9]{6,9}$`
	re := regexp.MustCompile(pattern)

	// Generate a Valor string
	result := function.Valor()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated Valor %s does not match the expected pattern", result)
	}
}

func TestWkn(t *testing.T) {
	// Define the expected pattern for Wkn
	pattern := `^[ABCDEFGHLMNPQRSTUVXYZ]{6}$`
	re := regexp.MustCompile(pattern)

	// Generate a Wkn string
	result := function.Wkn()

	// Check if the result matches the pattern
	if !re.MatchString(result) {
		t.Errorf("Generated Wkn %s does not match the expected pattern", result)
	}
}
