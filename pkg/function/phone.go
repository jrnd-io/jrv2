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

package function

import (
	"text/template"

	"github.com/biter777/countries"
	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/random"
)

func init() {
	AddFuncs(template.FuncMap{
		"country_code":    CountryCode,
		"country_code_at": CountryCodeAt,
		"imei":            Imei,
		"phone":           Phone,
		"phone_at":        PhoneAt,
		"mobile_phone":    MobilePhone,
		"mobile_phone_at": MobilePhoneAt,
	})
}

const (
	PhoneMap       = "phone"
	MobilePhoneMap = "mobile_phone"
)

// CountryCode returns a random Country Code prefix
func CountryCode() string {
	countryIndex := emitter.GetState().CountryIndex
	if countryIndex == -1 {
		index := random.Random.IntN(len(countries.All()))
		return countries.ByNumeric(index).Info().CallCodes[0].String()
	}

	return countries.ByNumeric(countryIndex).Info().CallCodes[0].String()
}

// CountryCodeAt returns a Country Code prefix at a given index
func CountryCodeAt(index int) string {
	return countries.ByNumeric(index).Info().CallCodes[0].String()
}

// Imei returns a random imei number of 15 digits
func Imei() string {
	account := make([]byte, 14)
	for i := range account {
		account[i] = digits[random.Random.IntN(len(digits))]
	}
	first14 := string(account)
	return first14 + LuhnCheckDigit(first14)
}

// Phone returns a random land prefix
func Phone() string {
	cityIndex := emitter.GetState().CityIndex
	if cityIndex == -1 {
		l := Word(PhoneMap)
		lp, _ := Regex(l)
		return lp
	}

	return PhoneAt(cityIndex)
}

// PhoneAt returns a land prefix at a given index
func PhoneAt(index int) string {
	l := WordAt(PhoneMap, index)
	lp, _ := Regex(l)
	return lp
}

// MobilePhone returns a random mobile phone
func MobilePhone() string {
	countryIndex := emitter.GetState().CountryIndex
	if countryIndex == -1 {
		m := Word(MobilePhoneMap)
		mp, _ := Regex(m)
		return mp
	}

	return MobilePhoneAt(countryIndex)
}

// MobilePhoneAt returns a mobile phone at a given index
func MobilePhoneAt(index int) string {
	m := WordAt(MobilePhoneMap, index)
	mp, _ := Regex(m)
	return mp
}
