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
	"strconv"
	"testing"

	"github.com/biter777/countries"
	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

const (
	TestLocale = "test"
	City       = "city"
	Capital    = "capital"
	Country    = "country"
)

func TestBuildingNumber(t *testing.T) {
	n := 10
	s := function.BuildingNumber(n)

	assert.LessOrEqual(t, len(s), n)
}

func TestCapital(t *testing.T) {
	emitter.GetState().Locale = TestLocale
	// preload cache
	function.ClearCache(Capital)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", Capital), Capital)
	if err != nil {
		t.Error(err)
		return
	}

	c := function.Capital()
	assert.Contains(t, function.GetCache(Capital), c)
	function.ClearCache(Capital)
}
func TestCapitalAt(t *testing.T) {
	emitter.GetState().Locale = TestLocale
	// preload cache
	function.ClearCache(Capital)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", Capital), Capital)
	if err != nil {
		t.Error(err)
		return
	}

	c := function.CapitalAt(2)
	assert.Equal(t, "capital02", c)
	function.ClearCache(Capital)
}

func TestCardinal(t *testing.T) {
	dShort := function.Cardinal(true)
	dLong := function.Cardinal(false)

	assert.Contains(t, function.CardinalShort, dShort)
	assert.Contains(t, function.CardinalLong, dLong)
}

func TestCity(t *testing.T) {
	emitter.GetState().Locale = TestLocale
	// preload cache
	function.ClearCache(City)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", City), City)
	if err != nil {
		t.Error(err)
		return
	}

	c := function.City()
	assert.Contains(t, function.GetCache(City), c)
	function.ClearCache(City)
}
func TestCityAt(t *testing.T) {
	emitter.GetState().Locale = TestLocale
	// preload cache
	function.ClearCache(City)
	_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", City), City)
	if err != nil {
		t.Error(err)
		return
	}

	c := function.CityAt(2)
	assert.Equal(t, "city02", c)
	function.ClearCache(City)
}

func TestCountry(t *testing.T) {
	emitter.GetState().Locale = TestLocale

	c := function.Country()
	country := countries.ByName(c)
	assert.Contains(t, country.Alpha2(), c)
}
func TestCountryAt(t *testing.T) {

	c := function.CountryAt(232)
	assert.Equal(t, countries.ByNumeric(232).Alpha2(), c)
}

func TestLatitude(t *testing.T) {
	l := function.Latitude()
	lf, err := strconv.ParseFloat(l, 64)
	if err != nil {
		t.Error(err)
	}
	assert.InDelta(t, lf, -90.0, 180.0)
}

func TestLongitude(t *testing.T) {
	l := function.Longitude()
	lf, err := strconv.ParseFloat(l, 64)
	if err != nil {
		t.Error(err)
	}
	assert.InDelta(t, lf, -180.0, 360.0)
}

func TestNearbyGPS(t *testing.T) {
	t.Error()
}

func TestState(t *testing.T) {
	t.Error()
}
func TestStateAt(t *testing.T) {
	t.Error()
}
func TestStateShort(t *testing.T) {
	t.Error()
}
func TestStateShortAt(t *testing.T) {
	t.Error()
}
func TestStreet(t *testing.T) {
	t.Error()
}
func TestZip(t *testing.T) {
	t.Error()
}
func TestZipAt(t *testing.T) {
	t.Error()
}
