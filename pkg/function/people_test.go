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
	"testing"

	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/squeeze69/generacodicefiscale"
	"github.com/stretchr/testify/assert"
)

func TestPeopleFun(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name    string
		f       funcT
		funcMap string
	}{
		{
			name:    function.CompanyMap,
			f:       function.Company,
			funcMap: function.CompanyMap,
		},
		{
			name:    function.MailProviderMap,
			f:       function.EmailProvider,
			funcMap: function.MailProviderMap,
		},
		{
			name:    function.SurnameMap,
			f:       function.Surname,
			funcMap: function.SurnameMap,
		},
		{
			name:    function.NameFMap,
			f:       function.NameF,
			funcMap: function.NameFMap,
		},
		{
			name:    function.NameMMap,
			f:       function.NameM,
			funcMap: function.NameMMap,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			emitter.GetState().Locale = TestLocale
			function.ClearCache(tc.funcMap)
			_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", tc.funcMap), tc.funcMap)
			if err != nil {
				t.Error(err)
				return
			}
			c := tc.f()
			assert.Contains(t, function.GetCache(tc.funcMap), c)
			function.ClearCache(tc.funcMap)

		})
	}
}

func TestSsn(t *testing.T) {
	// Call the function

	// Define the expected SSN format using a regular expression
	ssnPattern := `^\d{3}-\d{2}-\d{4}$`

	re := regexp.MustCompile(ssnPattern)
	for i := 0; i < 100; i++ {
		ssn := function.Ssn()
		matched := re.MatchString(ssn)

		assert.True(t, matched, "SSN does not match the expected format")
	}

}

func TestUsername(t *testing.T) {
	firstName := "John"
	lastName := "Doe"

	// Define the expected username patterns using regular expressions
	patterns := []string{
		"^j[ohn]*[.\\-_]*[doe]*$",
	}

	for i := 0; i < 100; i++ {
		username := function.Username(firstName, lastName)
		matched := false
		for _, pattern := range patterns {
			matched, _ = regexp.MatchString(pattern, username)
			if matched {
				break
			}
		}
		// Assert that the username matches one of the expected patterns
		assert.True(t, matched, fmt.Sprintf("Username %s does not match any of the expected patterns", username))
	}
}

func TestUser(t *testing.T) {
	firstName := "John"
	lastName := "Doe"
	size := 8

	// Define the expected username pattern using a regular expression
	usernamePattern := `^[a-zA-Z]{1,8}\d{2}$`

	re := regexp.MustCompile(usernamePattern)

	for i := 0; i < 100; i++ {
		// Call the function
		username := function.User(firstName, lastName, size)

		// Check if the username matches the expected pattern
		matched := re.MatchString(username)

		// Assert that the username matches the expected pattern
		assert.True(t, matched, "Username does not match the expected pattern")
	}
}

func TestCodiceFiscale(t *testing.T) {
	testCases := []struct {
		name     string
		cityName string
		v        map[string]interface{}
	}{
		{
			name: "test_roma",
			v: map[string]interface{}{
				"_name":      "Mario",
				"_surname":   "Rossi",
				"_gender":    "M",
				"_birthdate": "1980-01-01",
				"_city":      "Roma",
			},
			cityName: "Roma",
		},
		{
			name: "test_bolzano",
			v: map[string]interface{}{
				"_name":      "Mario",
				"_surname":   "Bianchi",
				"_gender":    "M",
				"_birthdate": "1980-01-01",
				"_city":      "Bolzano",
			},
			cityName: "Bolzano/Bozen",
		},
		{
			name: "test_reggio",
			v: map[string]interface{}{
				"_name":      "Mario",
				"_surname":   "Bianchi",
				"_gender":    "M",
				"_birthdate": "1980-01-01",
				"_city":      "Reggio Emilia",
			},
			cityName: "Reggio nell'Emilia",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.v {
				emitter.GetState().Ctx.Store(k, v)
			}
			cf := function.CodiceFiscale()

			city, err := generacodicefiscale.CercaComune(tc.cityName)
			assert.Nil(t, err)
			expected, erg := generacodicefiscale.Genera(
				tc.v["_surname"].(string),
				tc.v["_name"].(string),
				tc.v["_gender"].(string),
				city.Codice,
				tc.v["_birthdate"].(string),
			)
			assert.Nil(t, erg)
			assert.NotEmpty(t, cf)
			assert.Equal(t, expected, cf)
		})
	}

}
