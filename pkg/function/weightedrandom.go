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
	"errors"
	"math/rand"
	"text/template"
)

func init() {
	AddFuncs(template.FuncMap{
		"weighted_string": WeightedRandomString,
		"weighted_int":    WeightedRandomInt,
	})

}

// WeightedRandomString selects a random string from the given slice
// with probability proportional to the corresponding weight.
func WeightedRandomString(items []string, weights []float64) (string, error) {
	if len(items) != len(weights) {
		return "", errors.New("items and weights slices must have the same length")
	}

	if len(items) == 0 {
		return "", errors.New("items slice cannot be empty")
	}

	// Calculate the sum of weights
	totalWeight := 0.0
	for _, w := range weights {
		if w < 0 {
			return "", errors.New("weights must be non-negative")
		}
		totalWeight += w
	}

	// Generate a random number between 0 and totalWeight
	r := rand.Float64() * totalWeight //nolint no need for a secure random generator

	// Find the selected item
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return items[i], nil
		}
	}

	// This should never happen, but return the last item just in case
	return items[len(items)-1], nil
}

// WeightedRandomInt selects a random integer from the given slice
// with probability proportional to the corresponding weight.
func WeightedRandomInt(items []int, weights []float64) (int, error) {
	if len(items) != len(weights) {
		return 0, errors.New("items and weights slices must have the same length")
	}

	if len(items) == 0 {
		return 0, errors.New("items slice cannot be empty")
	}

	// Calculate the sum of weights
	totalWeight := 0.0
	for _, w := range weights {
		if w < 0 {
			return 0, errors.New("weights must be non-negative")
		}
		totalWeight += w
	}

	// Generate a random number between 0 and totalWeight
	r := rand.Float64() * totalWeight //nolint no need for a secure random generator

	// Find the selected item
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return items[i], nil
		}
	}

	// This should never happen, but return the last item just in case
	return items[len(items)-1], nil
}
