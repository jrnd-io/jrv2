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
	"testing"
	"time"

	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

const (
	DateFrom   = "2023-01-01"
	DateTo     = "2023-12-31"
	DateFormat = "2006-01-02"
)

func TestUnixTimeStamp(t *testing.T) {
	// Test case when a positive number of days is provided
	days := 10
	now := time.Now()
	past := now.AddDate(0, 0, -days).Unix()
	result := function.UnixTimeStamp(days)
	assert.GreaterOrEqual(t, result, past, "Timestamp should be greater than or equal to the past timestamp")
	assert.LessOrEqual(t, result, now.Unix(), "Timestamp should be less than or equal to the current timestamp")

	// Test case when zero days are provided
	days = 0
	now = time.Now()
	past = now.AddDate(0, 0, -1).Unix()
	result = function.UnixTimeStamp(days)
	assert.GreaterOrEqual(t, result, past, "Timestamp should be greater than or equal to the past 24 hours timestamp")
	assert.LessOrEqual(t, result, now.Unix(), "Timestamp should be less than or equal to the current timestamp")
}

func TestDateBetween(t *testing.T) {
	// Test case when valid fromDate and toDate are provided
	fromDate := DateFrom
	toDate := DateTo
	result := function.DateBetween(fromDate, toDate)

	// Parse the result to ensure it's a valid date
	resultDate, err := time.Parse(time.DateOnly, result)
	assert.NoError(t, err, "Result should be a valid date")

	// Parse the fromDate and toDate for comparison
	startDate, err := time.Parse(time.DateOnly, fromDate)
	assert.NoError(t, err, "fromDate should be a valid date")
	endDate, err := time.Parse(time.DateOnly, toDate)
	assert.NoError(t, err, "toDate should be a valid date")

	// Ensure the result date is within the range
	assert.True(t, resultDate.After(startDate) || resultDate.Equal(startDate), "Result date should be after or equal to fromDate")
	assert.True(t, resultDate.Before(endDate) || resultDate.Equal(endDate), "Result date should be before or equal to toDate")

	fromDate = DateTo
	toDate = DateFrom
	result = function.DateBetween(fromDate, toDate)

	// Parse the result to ensure it's a valid date
	resultDate, err = time.Parse(time.DateOnly, result)
	assert.NoError(t, err, "Result should be a valid date")

	// Parse the fromDate and toDate for comparison
	startDate, err = time.Parse(time.DateOnly, fromDate)
	assert.NoError(t, err, "fromDate should be a valid date")
	endDate, err = time.Parse(time.DateOnly, toDate)
	assert.NoError(t, err, "toDate should be a valid date")

	// Ensure the result date is within the range
	assert.True(t, resultDate.After(endDate) || resultDate.Equal(endDate), "Result date should be after or equal to toDate")
	assert.True(t, resultDate.Before(startDate) || resultDate.Equal(startDate), "Result date should be before or equal to fromDate")
}

func TestDatesBetween(t *testing.T) {
	// Test case when valid fromDate, toDate, and num are provided
	fromDate := DateFrom
	toDate := DateTo
	num := 5
	results := function.DatesBetween(fromDate, toDate, num)

	// Ensure the length of the results matches num
	assert.Equal(t, num, len(results), "The length of the results should match num")

	// Parse the fromDate and toDate for comparison
	startDate, err := time.Parse(DateFormat, fromDate)
	assert.NoError(t, err, "fromDate should be a valid date")
	endDate, err := time.Parse(DateFormat, toDate)
	assert.NoError(t, err, "toDate should be a valid date")

	// Ensure each date in the results is within the range
	for _, result := range results {
		resultDate, err := time.Parse(DateFormat, result)
		assert.NoError(t, err, "Result should be a valid date")
		assert.True(t, resultDate.After(startDate) || resultDate.Equal(startDate), "Result date should be after or equal to fromDate")
		assert.True(t, resultDate.Before(endDate) || resultDate.Equal(endDate), "Result date should be before or equal to toDate")
	}
}

func TestBirthDate(t *testing.T) {
	minAge := 18
	maxAges := []int{18, 28, 38}

	for _, maxAge := range maxAges {
		birthDateStr := function.BirthDate(minAge, maxAge)
		birthDate, err := time.Parse(DateFormat, birthDateStr)
		assert.NoError(t, err, "BirthDate should return a valid date string")

		currentYear := time.Now().Year()
		maxBirthYear := currentYear - minAge
		minBirthYear := maxBirthYear - (maxAge - minAge)

		// Ensure the birth year is within the expected range
		assert.True(t, birthDate.Year() >= minBirthYear && birthDate.Year() <= maxBirthYear, "Birth year should be within the expected range")

		// Ensure the birth month is between 1 and 12
		assert.True(t, birthDate.Month() >= 1 && birthDate.Month() <= 12, "Birth month should be between 1 and 12")

		// Ensure the birth day is valid for the given month and year
		lastDayOfMonth := time.Date(birthDate.Year(), birthDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		assert.True(t, birthDate.Day() >= 1 && birthDate.Day() <= lastDayOfMonth, "Birth day should be valid for the given month and year")
	}
}

func TestPast(t *testing.T) {
	ayears := []int{0, 5, 10}

	for _, years := range ayears {
		pastDateStr := function.Past(years)
		pastDate, err := time.Parse(DateFormat, pastDateStr)
		assert.NoError(t, err, "Past should return a valid date string")

		now := time.Now().UTC()
		expectedStart := now.AddDate(-years, 0, 0)

		// Ensure the returned date is within the expected range
		if years == 0 {
			assert.True(t, pastDate.Before(now), "Past date should be within the expected range")
		} else {
			assert.True(t, pastDate.After(expectedStart) && pastDate.Before(now), "Past date should be within the expected range")
		}
	}
}

func TestFuture(t *testing.T) {
	ayears := []int{0, 5, 10}

	for _, years := range ayears {
		futureDateStr := function.Future(years)
		futureDate, err := time.Parse(DateFormat, futureDateStr)
		assert.NoError(t, err, "Future should return a valid date string")

		now := time.Now().UTC()
		expectedEnd := now.AddDate(years, 0, 0)

		// Ensure the returned date is within the expected range
		if years == 0 {
			assert.True(t, futureDate.Before(now), "Past date should be within the expected range")
		} else {
			assert.True(t, futureDate.After(now) && futureDate.Before(expectedEnd), "Future date should be within the expected range")
		}
	}
}

func TestRecent(t *testing.T) {
	adays := []int{0, 10}

	for _, days := range adays {
		recentDateStr := function.Recent(days)
		recentDate, err := time.Parse(DateFormat, recentDateStr)
		assert.NoError(t, err, "Recent should return a valid date string")

		now := time.Now().UTC()
		expectedStart := now.AddDate(0, 0, -days)

		// Ensure the returned date is within the expected range
		if days == 0 {
			assert.True(t, recentDate.Before(now), "Recent date should be before now")
		} else {
			assert.True(t, recentDate.After(expectedStart) && recentDate.Before(now), "Recent date should be within the expected range")

		}
	}
}

func TestSoon(t *testing.T) {
	adays := []int{0, 10}

	for _, days := range adays {
		soonDateStr := function.Soon(days)
		soonDate, err := time.Parse(DateFormat, soonDateStr)
		assert.NoError(t, err, "Soon should return a valid date string")

		now := time.Now().UTC()
		expectedEnd := now.AddDate(0, 0, days)

		// Ensure the returned date is within the expected range
		if days == 0 {
			assert.True(t, soonDate.Before(now), "Recent date should be before now")
		} else {
			assert.True(t, soonDate.After(now) && soonDate.Before(expectedEnd), "Soon date should be within the expected range")
		}
	}
}
