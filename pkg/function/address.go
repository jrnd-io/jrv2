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
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/config"
	"math"
	"text/template"

	"github.com/biter777/countries"
	"github.com/jrnd-io/jrv2/pkg/emitter"
)

func init() {
	AddFuncs(template.FuncMap{
		"building":       BuildingNumber,
		"cardinal":       Cardinal,
		"capital":        Capital,
		"capital_at":     CapitalAt,
		"city":           City,
		"city_at":        CityAt,
		"country":        Country,
		"country_random": CountryRandom,
		"country_at":     CountryAt,
		"latitude":       Latitude,
		"longitude":      Longitude,
		"nearby_gps":     NearbyGPS,
		"state":          State,
		"state_at":       StateAt,
		"state_short":    StateShort,
		"state_short_at": StateShortAt,
		"street":         Street,
		"zip":            Zip,
		"zip_at":         ZipAt,
	})
}

const (
	earthRadius     = 6371000 // in meters
	degreesPerMeter = 1.0 / earthRadius * 180.0 / math.Pi

	CityMap       = "city"
	CapitalMap    = "capital"
	StateMap      = "state"
	StateShortMap = "state"
	StreetMap     = "street"
	ZipMap        = "zip"
)

var CardinalShort = []string{"N", "S", "E", "O", "NE", "NO", "SE", "SO"}
var CardinalLong = []string{"North", "South", "East", "Ovest", "North-East", "North-Ovest", "South-East", "South-Ovest"}

// BuildingNumber generates a random building number of max n digits
func BuildingNumber(n int) string {
	building := make([]byte, config.Random.Intn(n)+1)
	for i := range building {
		building[i] = digits[config.Random.Intn(len(digits))]
	}
	return string(building)
}

// Capital returns a random Capital
func Capital() string {
	return Word(CapitalMap)
}

// CapitalAt returns Capital at given index
func CapitalAt(index int) string {
	return WordAt(CapitalMap, index)
}

// Cardinal return a random cardinal direction, in long or short form
func Cardinal(short bool) string {

	directions := CardinalLong
	if short {
		directions = CardinalShort
	}

	return directions[config.Random.Intn(len(directions))]
}

// City returns a random City
func City() string {
	c := Word(CityMap)
	emitter.GetState().Ctx.Store(fmt.Sprintf("_%s", CityMap), c)
	emitter.GetState().CityIndex = emitter.GetState().LastIndex
	return c
}

// CityAt returns City at given index
func CityAt(index int) string {
	return WordAt(CityMap, index)
}

// Country returns the ISO 3166 Country selected with locale
func Country() string {
	countryIndex := emitter.GetState().CountryIndex
	if countryIndex == -1 || countryIndex == 0 {
		emitter.GetState().LastIndex = config.Random.Intn(len(countries.All()))
		c := countries.All()[emitter.GetState().LastIndex].Alpha2()
		return c
	}

	return countries.ByNumeric(emitter.GetState().CountryIndex).Alpha2()
}

// CountryRandom returns a random ISO 3166 Country
func CountryRandom() string {
	emitter.GetState().LastIndex = config.Random.Intn(len(countries.All()))
	return countries.ByNumeric(emitter.GetState().LastIndex).Alpha2()
}

// CountryAt returns an ISO 3166 Country at a given index
func CountryAt(index int) string {
	return countries.ByNumeric(index).Alpha2()
}

// Latitude returns a random latitude between -90 and 90
func Latitude() string {
	latitude := -90 + config.Random.Float64()*(180)
	return fmt.Sprintf("%.4f", latitude)
}

// Longitude returns a random longitude between -180 and 180
func Longitude() string {
	longitude := -180 + config.Random.Float64()*(360)
	return fmt.Sprintf("%.4f", longitude)
}

// NearbyGPS returns a random latitude longitude within a given radius in meters
func NearbyGPS(latitude float64, longitude float64, radius int) string {
	radiusInMeters := float64(radius)

	// Generate a random angle in radians
	randomAngle := config.Random.Float64() * 2 * math.Pi

	// Calculate the distance from the center point
	distanceInMeters := config.Random.Float64() * radiusInMeters

	// Convert the distance to degrees
	distanceInDegrees := distanceInMeters * degreesPerMeter

	// Calculate the new latitude and longitude
	newLatitude := latitude + (distanceInDegrees * math.Cos(randomAngle))
	newLongitude := longitude + (distanceInDegrees * math.Sin(randomAngle))

	return fmt.Sprintf("%.4f %.4f", newLatitude, newLongitude)

}

// State returns a random State
func State() string {
	s := Word(StateMap)
	emitter.GetState().Ctx.Store(fmt.Sprintf("_%s", StateMap), s)
	emitter.GetState().CountryIndex = emitter.GetState().LastIndex
	return s
}

// StateAt returns State at given index
func StateAt(index int) string {
	return WordAt(StateMap, index)
}

// StateShort returns a random short State
func StateShort() string {
	return Word(StateShortMap)
}

// StateShortAt returns short State at given index
func StateShortAt(index int) string {
	return WordAt(StateShortMap, index)
}

// Street returns a random street
func Street() string {
	return Word(StreetMap)
}

// StreetAt returns a street at given index
func StreetAt(index int) string {
	return WordAt(StreetMap, index)
}

// Zip returns a random Zip code
func Zip() string {
	cityIndex := emitter.GetState().CityIndex

	if cityIndex == -1 {
		z := Word(ZipMap)
		zip, _ := Regex(z)
		return zip
	}

	return ZipAt(cityIndex)
}

// ZipAt returns Zip code at given index
func ZipAt(index int) string {
	z := WordAt(ZipMap, index)
	zip, _ := Regex(z)
	return zip
}
