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

// modified from https://golang.org/doc%2Fcodewalk%2Fmarkov.go

package function

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/config"
	"io"
	"strings"
	"text/template"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

//go:embed alice.txt
var alice string

//go:embed loremipsum.txt
var lorem string

func init() {
	AddFuncs(template.FuncMap{
		"sentence":        Sentence,       // to nonsense
		"sentence_prefix": SentencePrefix, // to nonsense
		"lorem":           Lorem,          // to nonsense
		"markov":          Nonsense,       // to nonsense
	})

}

// Get Alice text
func GetAlice() string {
	return alice
}

// Get Lorem text
func GetLorem() string {
	return lorem
}

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[config.Random.Intn(len(choices))]

		if i == n-1 {
			if strings.HasSuffix(next, ",") {
				next = strings.ReplaceAll(next, ",", ".")
			}
		}
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

// Lorem generates a 'lorem ipsum' text of size words
func Lorem(size int) string {
	return Nonsense(2, size, string(lorem))
}

// SentencePrefix generates an 'alice in wonderland' text of size words with given prefixLen
func SentencePrefix(prefixLen, numWords int) string {
	return Nonsense(prefixLen, numWords, string(alice))
}

// Nonsense generates a random Sentence of numWords wordsm using a prefixLen and a baseText to start from
func Nonsense(prefixLen, numWords int, baseText string) string {
	c := NewChain(prefixLen)
	c.Build(strings.NewReader(baseText))
	return c.Generate(numWords)
}

// RandomStringVocabulary returns a random string long between min and max characters using a vocabulary
func RandomStringVocabulary(min, max int, source string) string {
	if len(source) == 0 {
		return ""
	}
	textb := make([]byte, min+config.Random.Intn(max-min))
	for i := range textb {
		textb[i] = source[config.Random.Intn(len(source))]
	}
	return string(textb)
}

// Sentence generates an 'alice in wonderland' text of size words
func Sentence(numWords int) string {
	return SentencePrefix(2, numWords)
}
