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
)

func init() {
	AddFuncs(template.FuncMap{
		"add":          func(a, b int) int { return a + b },
		"div":          func(a, b int) int { return a / b },
		"format_float": func(f string, v float32) string { return fmt.Sprintf(f, v) },
		"integer":      func(min, max int) int { return min + config.Random.Intn(max-min) },
		"integer64":    func(min, max int64) int64 { return min + config.Random.Int63n(max-min) },
		"floating":     func(min, max float32) float32 { return min + config.Random.Float32()*(max-min) },
		"sub":          func(a, b int) int { return a - b },
		"max":          math.Max,
		"min":          math.Min,
		"minint":       Minint,
		"maxint":       Maxint,
		"mod":          func(a, b int) int { return a % b },
		"mul":          func(a, b int) int { return a * b },
	})

}

// Minint returns the minimum between two ints
func Minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Maxint returns the minimum between two ints
func Maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}
