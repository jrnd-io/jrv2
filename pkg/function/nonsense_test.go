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
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

func TestShift(t *testing.T) {
	tests := []struct {
		initial  function.Prefix
		word     string
		expected function.Prefix
	}{
		{function.Prefix{"a", "b", "c"}, "d", function.Prefix{"b", "c", "d"}},
		{function.Prefix{"x", "y", "z"}, "w", function.Prefix{"y", "z", "w"}},
	}

	for _, tt := range tests {
		tt.initial.Shift(tt.word)
		if !reflect.DeepEqual(tt.initial, tt.expected) {
			t.Errorf("Shift(%q) = %v; want %v", tt.word, tt.initial, tt.expected)
		}
	}
}

func TestLorem(t *testing.T) {
	t.Skip("This test is flaky")
	tests := []struct {
		name     string
		size     int
		expected int // Expected length of the generated string
	}{
		{
			name:     "small size",
			size:     10,
			expected: 10,
		},
		{
			name:     "medium size",
			size:     50,
			expected: 50,
		},
		{
			name:     "large size",
			size:     100,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := function.Lorem(tt.size)
			if len(strings.Fields(result)) > tt.expected {
				t.Errorf("Lorem(%d) = %d words; want %d words", tt.size, len(strings.Fields(result)), tt.expected)
			}
		})
	}
}
func TestSentencePrefix(t *testing.T) {
	tests := []struct {
		name      string
		prefixLen int
		numWords  int
	}{
		{
			name:      "short sentence",
			prefixLen: 1,
			numWords:  5,
		},
		{
			name:      "medium sentence",
			prefixLen: 2,
			numWords:  10,
		},
		{
			name:      "long sentence",
			prefixLen: 3,
			numWords:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := function.SentencePrefix(tt.prefixLen, tt.numWords)
			words := strings.Fields(result)
			if len(words) != tt.numWords {
				t.Errorf("SentencePrefix(%d, %d) = %d words; want %d words", tt.prefixLen, tt.numWords, len(words), tt.numWords)
			}
		})
	}
}

func TestRandomStringVocabulary(t *testing.T) {
	tests := []struct {
		name   string
		min    int
		max    int
		source string
	}{
		{
			name:   "basic test",
			min:    5,
			max:    10,
			source: "abc",
		},
		{
			name:   "single character source",
			min:    3,
			max:    6,
			source: "x",
		},
		{
			name:   "empty source",
			min:    0,
			max:    5,
			source: "",
		},
		{
			name:   "large range",
			min:    10,
			max:    100,
			source: "abcdefghijklmnopqrstuvwxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := function.RandomStringVocabulary(tt.min, tt.max, tt.source)
			length := len(result)
			if length < tt.min || length > tt.max {
				t.Errorf("RandomStringVocabulary(%d, %d, %q) = %d characters; want length between %d and %d", tt.min, tt.max, tt.source, length, tt.min, tt.max)
			}
			for _, char := range result {
				if !strings.ContainsRune(tt.source, char) {
					t.Errorf("RandomStringVocabulary(%d, %d, %q) generated character %q not in source %q", tt.min, tt.max, tt.source, char, tt.source)
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		prefixLen int
		n         int
		expected  string
	}{
		{
			name:      "basic test",
			input:     "hello world",
			prefixLen: 1,
			n:         2,
			expected:  "hello world",
		},
		{
			name:      "test with comma",
			input:     "hello world,",
			prefixLen: 1,
			n:         2,
			expected:  "hello world.",
		},
		{
			name:      "test with break",
			input:     "hello",
			prefixLen: 1,
			n:         2,
			expected:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := function.NewChain(tt.prefixLen)
			c.Build(bytes.NewBufferString(tt.input))
			result := c.Generate(tt.n)
			if result != tt.expected {
				t.Errorf("Generate(%d) = %q; want %q", tt.n, result, tt.expected)
			}
		})
	}
}

func TestSentence(t *testing.T) {
	tests := []struct {
		numWords int
		expected string
	}{
		{numWords: 5},
		{numWords: 10},
		{numWords: 0, expected: ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("numWords=%d", tt.numWords), func(t *testing.T) {
			result := function.Sentence(tt.numWords)
			assert.Equal(t, len(strings.Fields(result)), tt.numWords)
		})
	}
}
