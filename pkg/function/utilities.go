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
	"sort"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

func init() {
	AddFuncs(template.FuncMap{
		"array":    func(count int) []int { return make([]int, count) },
		"bool":     RandomBool,
		"image":    Image,
		"image_of": ImageOf,
		"index_of": IndexOf,
		"inject":   Inject,
		"key":      func(name string, n int) string { return fmt.Sprintf("%s%d", name, config.Random.IntN(n)) },
		"seed":     Seed,
		"uuid":     UniqueID,
		"xid":      Xid,
		"yesorno":  YesOrNo,
	})

}

// Image generates a random Image url of given width, height and type
func Image(width int, height int) string {
	imageType := []string{"abstract", "animals", "business", "cats", "city", "fashion", "food", "nature", "nightlife", "people", "sport", "technics", "transport"}
	return ImageOf(
		width,
		height,
		imageType[config.Random.IntN(len(imageType))],
	)
}

// ImageOf generates a random Image url of given width, height and type
func ImageOf(width int, height int, imageType string) string {
	return fmt.Sprintf("https://loremflickr.com/%d/%d/%s", width, height, imageType)
}

// RandomBool returns a random boolean
func RandomBool() string {
	b := config.Random.IntN(2)
	if b == 0 {
		return "false"
	}

	return "true"
}

// UniqueId returns a random uuid
func UniqueID() string {
	return uuid.New().String()
}

// Xid returns a random uuid (xid is a globally unique id generator thought for the web)
func Xid() string { return xid.New().String() }

// YesOrNo returns a random yes or no
func YesOrNo() string {
	b := config.Random.IntN(2)
	if b == 0 {
		return "no"
	}

	return "yes"
}

// Contains checks if the str string is present in an array of string []string
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Inject is used to inject a different value with a given probability, typically used to generate a bad value
func Inject(probability float64, injected, original any) any {
	if config.Random.Float64() < probability {
		return injected
	}
	return original
}

// Seed sets seeds and can be used in a template
func Seed(rndSeed uint64) string {
	SetSeed(rndSeed)
	return ""
}

// SetSeed sets seeds for all random JR objects
func SetSeed(rndSeed uint64) {
	config.ChaCha8.Seed(config.CreateByteSeed(rndSeed))
	uuid.SetRand(config.ChaCha8)
}

// IndexOf returns the index of the s string in a file
func IndexOf(s string, name string) int {
	_, err := Cache(name)
	if err != nil {
		return -1
	}
	words := data[name]
	index := sort.Search(len(words), func(i int) bool { return strings.ToLower(words[i]) >= strings.ToLower(s) })

	if index < len(words) && words[index] == s {
		return index
	}

	return -1
}
