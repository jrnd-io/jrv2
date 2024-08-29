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
	"regexp"
	"testing"

	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/stretchr/testify/assert"
)

func TestRegex(t *testing.T) {
	// Test case when a valid regex pattern is provided
	validRegex := `^[a-z]{5}$`
	result, err := function.Regex(validRegex)
	assert.NoError(t, err, "Expected no error for valid regex")
	assert.Regexp(t, regexp.MustCompile(validRegex), result, "Result should match the regex pattern")

	// Test case when an invalid regex pattern is provided
	invalidRegex := `[a-z`
	_, err = function.Regex(invalidRegex)
	assert.Error(t, err, "Expected an error for invalid regex")

	// Test case for a complex regex pattern and a matching string
	complexRegex := `^([A-Za-z0-9]+[-_])*[A-Za-z0-9]+$`
	result, err = function.Regex(complexRegex)
	assert.NoError(t, err, "Expected no error for valid regex")
	assert.Regexp(t, regexp.MustCompile(complexRegex), result, "Matching string should match the regex pattern")

}
