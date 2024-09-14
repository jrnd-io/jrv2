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
	"math"
	"strings"
	"testing"

	"github.com/jrnd-io/jrv2/pkg/function"
)

func TestRandoms(t *testing.T) {
	// Test case 1: Single argument (normal random)
	t.Run("Single argument", func(t *testing.T) {
		input := "apple|banana|cherry"
		result := function.Randoms(input)
		if !strings.Contains(input, result) {
			t.Errorf("Expected result to be one of %s, but got %s", input, result)
		}
	})

	// Test case 2: Two arguments (weighted random)
	t.Run("Two arguments (weighted random)", func(t *testing.T) {
		input := "red|green|blue"
		weights := "0.5|0.3|0.2"

		// Run multiple times to check distribution
		counts := make(map[string]int)
		iterations := 1000

		for i := 0; i < iterations; i++ {
			result := function.Randoms(input, weights)
			counts[result]++
			if !strings.Contains(input, result) {
				t.Errorf("Expected result to be one of %s, but got %s", input, result)
			}
		}

		// Check if the distribution roughly matches the weights
		expectedRatios := map[string]float64{
			"red":   0.5,
			"green": 0.3,
			"blue":  0.2,
		}

		tolerance := 0.05 // Allow 5% tolerance
		for color, count := range counts {
			actualRatio := float64(count) / float64(iterations)
			expectedRatio := expectedRatios[color]
			if actualRatio < expectedRatio-tolerance || actualRatio > expectedRatio+tolerance {
				t.Errorf("Distribution for %s is off. Expected around %.2f, got %.2f", color, expectedRatio, actualRatio)
			}
		}
	})

	// Test case 3: Invalid weight
	t.Run("Invalid weight", func(t *testing.T) {
		input := "apple|banana"
		weights := "0.5|invalid"
		result := function.Randoms(input, weights)
		if result != "" {
			t.Errorf("Expected empty string for invalid weight, but got %s", result)
		}
	})
}

func TestWeightedRandomString(t *testing.T) {
	tests := []struct {
		name        string
		items       []string
		weights     []float64
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid input",
			items:   []string{"a", "b", "c"},
			weights: []float64{1.0, 2.0, 3.0},
			wantErr: false,
		},
		{
			name:        "Mismatched lengths",
			items:       []string{"a", "b"},
			weights:     []float64{1.0, 2.0, 3.0},
			wantErr:     true,
			errContains: "same length",
		},
		{
			name:        "Empty items",
			items:       []string{},
			weights:     []float64{},
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "Negative weight",
			items:       []string{"a", "b"},
			weights:     []float64{1.0, -1.0},
			wantErr:     true,
			errContains: "non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := function.WeightedRandomString(tt.items, tt.weights)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WeightedRandomString() error = nil, wantErr %v", tt.wantErr)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("WeightedRandomString() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("WeightedRandomString() unexpected error = %v", err)
				}
				if !containsString(tt.items, result) {
					t.Errorf("WeightedRandomString() result %v not in items %v", result, tt.items)
				}
			}
		})
	}
}

func TestWeightedRandomInt(t *testing.T) {
	tests := []struct {
		name        string
		items       []int
		weights     []float64
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid input",
			items:   []int{1, 2, 3},
			weights: []float64{1.0, 2.0, 3.0},
			wantErr: false,
		},
		{
			name:        "Mismatched lengths",
			items:       []int{1, 2},
			weights:     []float64{1.0, 2.0, 3.0},
			wantErr:     true,
			errContains: "same length",
		},
		{
			name:        "Empty items",
			items:       []int{},
			weights:     []float64{},
			wantErr:     true,
			errContains: "cannot be empty",
		},
		{
			name:        "Negative weight",
			items:       []int{1, 2},
			weights:     []float64{1.0, -1.0},
			wantErr:     true,
			errContains: "non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := function.WeightedRandomInt(tt.items, tt.weights)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WeightedRandomInt() error = nil, wantErr %v", tt.wantErr)
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("WeightedRandomInt() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("WeightedRandomInt() unexpected error = %v", err)
				}
				if !containsInt(tt.items, result) {
					t.Errorf("WeightedRandomInt() result %v not in items %v", result, tt.items)
				}
			}
		})
	}
}

func TestWeightedRandomDistribution(t *testing.T) {
	items := []string{"a", "b", "c"}
	weights := []float64{0.1, 0.2, 0.7}
	totalWeight := 1.0
	iterations := 100000

	counts := make(map[string]int)
	for i := 0; i < iterations; i++ {
		result, err := function.WeightedRandomString(items, weights)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		counts[result]++
	}

	for i, item := range items {
		expected := float64(iterations) * weights[i] / totalWeight
		actual := float64(counts[item])
		diff := math.Abs(actual - expected)
		tolerance := 0.05 * expected // 5% tolerance

		if diff > tolerance {
			t.Errorf("Distribution for item %s: expected around %.2f, got %.2f", item, expected, actual)
		}
	}
}

func TestWeightedRandomIntDistribution(t *testing.T) {
	items := []int{1, 2, 3}
	weights := []float64{0.1, 0.2, 0.7}
	totalWeight := 1.0
	iterations := 100000

	counts := make(map[int]int)
	for i := 0; i < iterations; i++ {
		result, err := function.WeightedRandomInt(items, weights)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		counts[result]++
	}

	for i, item := range items {
		expected := float64(iterations) * weights[i] / totalWeight
		actual := float64(counts[item])
		diff := math.Abs(actual - expected)
		tolerance := 0.05 * expected // 5% tolerance

		if diff > tolerance {
			t.Errorf("Distribution for item %d: expected around %.2f, got %.2f", item, expected, actual)
		}
	}
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsInt(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
