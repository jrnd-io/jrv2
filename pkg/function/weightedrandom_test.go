package function_test

import (
	"math"
	"strings"
	"testing"

	"github.com/jrnd-io/jrv2/pkg/function"
)

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
