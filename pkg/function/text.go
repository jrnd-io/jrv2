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
	"strconv"
	"strings"
	"text/template"

	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/random"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	AddFuncs(template.FuncMap{
		// text utilities
		"atoi":                     Atoi,
		"itoa":                     strconv.Itoa,
		"concat":                   func(a string, b string) string { return a + b },
		"counter":                  Counter,
		"first":                    func(s string) string { return s[:1] },
		"firstword":                func(s string) string { return strings.Split(s, " ")[0] },
		"from":                     Word,
		"from_at":                  WordAt,
		"from_shuffle":             WordShuffle,
		"from_n":                   WordShuffleN,
		"join":                     strings.Join,
		"len":                      Len,
		"lower":                    strings.ToLower,
		"random":                   func(s []string) string { return s[random.Random.IntN(len(s))] },
		"randoms":                  func(s string) string { a := strings.Split(s, "|"); return a[random.Random.IntN(len(a))] },
		"random_index":             RandomIndex,
		"random_string":            RandomString,
		"random_string_vocabulary": RandomStringVocabulary,
		"regex":                    Regex, // to regex
		"repeat":                   strings.Repeat,
		"replaceall":               strings.ReplaceAll,
		"squeeze":                  func(s string) string { return strings.ReplaceAll(s, " ", "") },
		"squeezechars":             func(s, c string) string { return strings.ReplaceAll(s, c, "") },
		"split":                    strings.Split,
		"substr":                   func(start, length int, s string) string { return s[start:length] },
		"trim":                     strings.TrimSpace,
		"trimchars":                strings.Trim,
		"title":                    cases.Title(language.English).String,
		"upper":                    strings.ToUpper,
	})

}

func Randoms(s ...string) string {
	a := strings.Split(s[0], "|")
	// use normal random if only one argument is provided
	if len(s) == 1 {
		a := strings.Split(s[0], "|")
		return a[random.Random.IntN(len(a))]
	}

	// if more than one argument is provided, second argument is a list of float separated by "|"
	// return a weighted random
	ws := strings.Split(s[1], "|")
	w := make([]float64, len(ws))
	for i := range ws {
		_w, err := strconv.ParseFloat(ws[i], 64)
		if err != nil {
			return ""
		}
		w[i] = _w
	}

	ret, err := WeightedRandomString(a, w)
	if err != nil {
		return ""
	}
	return ret

}

func Atoi(s string) int {
	if len(s) == 0 {
		return 0
	}

	i, _ := strconv.Atoi(s)
	return i
}

// Counter creates a counter named c, starting from start and incrementing by step
func Counter(c string, start, step int) int {
	return emitter.GetState().Counter(c, start, step)
}

// Word returns a random string from a list of strings in a file.
func Word(name string) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	emitter.GetState().LastIndex = random.Random.IntN(len(words))
	return words[emitter.GetState().LastIndex]
}

// WordAt returns a string at a given position in a list of strings in a file.
func WordAt(name string, index int) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	return words[index]
}

// WordShuffle returns a shuffled list of strings in a file.
func WordShuffle(name string) []string {
	_, err := Cache(name)
	if err != nil {
		return []string{""}
	}
	words := data[name]
	return WordShuffleN(name, len(words))
}

// wordShuffleN return a subset of n elements in a list of string in a file.
func WordShuffleN(name string, n int) []string {
	_, err := Cache(name)
	if err != nil {
		return []string{""}
	}
	words := data[name]
	random.Random.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	number := Minint(n, len(words))
	return words[:number]
}

// Len returns number of words (lines) in a word file
func Len(name string) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	l := len(data[name])
	return strconv.Itoa(l)
}

// RandomIndex returns a random index in a word file
func RandomIndex(name string) string {
	_, err := Cache(name)
	if err != nil {
		return ""
	}
	words := data[name]
	emitter.GetState().LastIndex = random.Random.IntN(len(words))
	return strconv.Itoa(emitter.GetState().LastIndex)
}

// RandomString returns a random string long between min and max characters
func RandomString(min, max int) string {
	return RandomStringVocabulary(min, max, alphabet)
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
	totalWeight := float64(0.0)
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
