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
	"github.com/jrnd-io/jrv2/pkg/random"
	"github.com/jrnd-io/jrv2/pkg/state"
	"math"
	"strconv"
	"testing"

	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

const (
	TestLocale = "test"
)

func TestBuildingNumber(t *testing.T) {
	n := 10
	s := function.BuildingNumber(n)

	assert.LessOrEqual(t, len(s), n)
}

type funcT func() string
type funcTAt func(int) string

func TestFun(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name    string
		f       funcT
		funcMap string
	}{
		{
			name:    function.CityMap,
			f:       function.City,
			funcMap: function.CityMap,
		},
		{
			name:    function.CapitalMap,
			f:       function.Capital,
			funcMap: function.CapitalMap,
		},
		{
			name:    function.StateMap,
			f:       function.State,
			funcMap: function.StateMap,
		},
		{
			name:    function.StateShortMap,
			f:       function.StateShort,
			funcMap: function.StateShortMap,
		},
		{
			name:    function.StreetMap,
			f:       function.Street,
			funcMap: function.StreetMap,
		},
		{
			name:    function.ZipMap,
			f:       function.Zip,
			funcMap: function.ZipMap,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			state.GetSharedState().Locale = TestLocale
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

func TestFunAt(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		f        funcTAt
		funcMap  string
		index    int
		expected string
	}{
		{
			name:     function.CityMap,
			f:        function.CityAt,
			funcMap:  function.CityMap,
			index:    2,
			expected: fmt.Sprintf("%s02", function.CityMap),
		},
		{
			name:     function.CapitalMap,
			f:        function.CapitalAt,
			funcMap:  function.CapitalMap,
			index:    2,
			expected: fmt.Sprintf("%s02", function.CapitalMap),
		},
		{
			name:     function.StateMap,
			f:        function.StateAt,
			funcMap:  function.StateMap,
			index:    2,
			expected: fmt.Sprintf("%s02", function.StateMap),
		},
		{
			name:     function.StateShortMap,
			f:        function.StateShortAt,
			funcMap:  function.StateShortMap,
			index:    2,
			expected: fmt.Sprintf("%s02", function.StateShortMap),
		},
		{
			name:     function.StreetMap,
			f:        function.StreetAt,
			funcMap:  function.StreetMap,
			index:    2,
			expected: fmt.Sprintf("%s02", function.StreetMap),
		},
		{
			name:     function.ZipMap,
			f:        function.ZipAt,
			funcMap:  function.ZipMap,
			index:    2,
			expected: "00002",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			state.GetSharedState().Locale = TestLocale
			function.ClearCache(tc.funcMap)
			_, err := function.CacheFromFile(fmt.Sprintf("./testdata/%s.txt", tc.funcMap), tc.funcMap)
			if err != nil {
				t.Error(err)
				return
			}
			c := tc.f(tc.index)
			assert.Equal(t, tc.expected, c)
			function.ClearCache(tc.funcMap)

		})
	}
}

func TestCardinal(t *testing.T) {
	dShort := function.Cardinal(true)
	dLong := function.Cardinal(false)

	assert.Contains(t, function.CardinalShort, dShort)
	assert.Contains(t, function.CardinalLong, dLong)
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
	// Define test cases
	testCases := []struct {
		latitude  float64
		longitude float64
		radius    int
	}{
		{37.7749, -122.4194, 1000}, // San Francisco, 1000 meters
		{51.5074, -0.1278, 500},    // London, 500 meters
		{-33.8688, 151.2093, 2000}, // Sydney, 2000 meters
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Lat: %f, Lon: %f, Radius: %d", tc.latitude, tc.longitude, tc.radius), func(t *testing.T) {
			result := function.NearbyGPS(tc.latitude, tc.longitude, tc.radius)
			var newLat, newLon float64
			np, err := fmt.Sscanf(result, "%f %f", &newLat, &newLon)
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, 2, np)

			// Calculate the distance between the original and new coordinates
			distance := haversine(tc.latitude, tc.longitude, newLat, newLon)

			// Check if the distance is within the radius
			if distance > float64(tc.radius) {
				t.Errorf("Expected distance <= %d meters, but got %f meters", tc.radius, distance)
			}
		})
	}
}

// haversine calculates the distance between two points on the Earth
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func init() {
	random.SetRandom(0)
}
