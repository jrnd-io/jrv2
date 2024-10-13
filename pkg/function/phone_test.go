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
	"github.com/jrnd-io/jrv2/pkg/state"
	"strconv"
	"testing"
	"unicode"

	"github.com/biter777/countries"
	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

func TestImei(t *testing.T) {
	imei := function.Imei()

	// Check that the length of the IMEI is 15
	assert.Equal(t, 15, len(imei), "IMEI length should be 15")

	// Check that the IMEI contains only digits
	for _, char := range imei {
		assert.True(t, unicode.IsDigit(char), "IMEI should contain only digits")
	}

	// Optionally, validate the IMEI using the Luhn algorithm
	assert.True(t, isValidLuhn(imei), "IMEI should be valid according to the Luhn algorithm")
}

// isValidLuhn validates a string using the Luhn algorithm
func isValidLuhn(s string) bool {
	var sum int
	nDigits := len(s)
	parity := nDigits % 2
	for i := 0; i < nDigits; i++ {
		digit := int(s[i] - '0')
		if i%2 == parity {
			digit *= 2
		}
		if digit > 9 {
			digit -= 9
		}
		sum += digit
	}
	return sum%10 == 0
}

func TestCountryCode(t *testing.T) {
	// Backup the original state
	// Test case when countryIndex is -1
	originalIndex := state.GetState().CountryIndex
	defer func() { state.GetState().CountryIndex = originalIndex }()

	state.GetState().CountryIndex = -1
	countryCode := function.CountryCode()
	assert.NotEmpty(t, countryCode, "Country code should not be empty when countryIndex is -1")

	index, err := strconv.Atoi(fmt.Sprintf("%d", countries.USA))
	assert.NoError(t, err, "Error should be nil when converting country code to integer")
	state.GetState().CountryIndex = index

	expectedCountryCode := countries.USA.Info().CallCodes[0].String()
	countryCode = function.CountryCode()
	assert.Equal(t, expectedCountryCode, countryCode, "Country code should match the expected value when countryIndex is a valid index")
}

func TestPhone(t *testing.T) {
	origCityIndex := state.GetState().CityIndex
	defer func() { state.GetState().CityIndex = origCityIndex }()

	state.GetState().Locale = TestLocale
	function.ClearCache(function.PhoneMap)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", function.PhoneMap), function.PhoneMap)
	assert.NoError(t, err, "Error should be nil when caching phone numbers")

	// Test case when cityIndex is -1
	state.GetState().CityIndex = -1
	phone := function.Phone()
	assert.NotEmpty(t, phone, "Phone number should not be empty when cityIndex is -1")

	// Test case when cityIndex is a valid index
	state.GetState().CityIndex = -1
	phone = function.PhoneAt(1)
	assert.NotEmpty(t, phone, "Phone number should not be empty when cityIndex is -1")
}

func TestMobilePhone(t *testing.T) {
	origIndex := state.GetState().CountryIndex
	defer func() { state.GetState().CountryIndex = origIndex }()

	state.GetState().Locale = TestLocale
	function.ClearCache(function.MobilePhoneMap)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", function.MobilePhoneMap), function.MobilePhoneMap)
	assert.NoError(t, err, "Error should be nil when caching phone numbers")

	// Test case when cityIndex is -1
	state.GetState().CountryIndex = -1
	phone := function.MobilePhone()
	assert.NotEmpty(t, phone, "Mobile Phone number should not be empty when countryIndex is -1")

	// Test case when cityIndex is a valid index
	state.GetState().CountryIndex = -1
	phone = function.MobilePhoneAt(1)
	assert.NotEmpty(t, phone, "Phone number should not be empty when countryIndex is -1")

}
