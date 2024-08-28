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
	"text/template"

	"github.com/jrnd-io/jrv2/pkg/emitter"
)

func init() {
	AddFuncs(template.FuncMap{
		// context utilities
		"add_v_to_list":            AddValueToList,
		"random_v_from_list":       RandomValueFromList,
		"random_n_v_from_list":     RandomNValuesFromList,
		"get_v_from_list_at_index": GetValueFromListAtIndex,
		"get_v":                    GetV,
		"set_v":                    SetV,
		"fromcsv":                  FromCSV,
	})

}

// AddValueToList adds value v to Context list l
func AddValueToList(l string, v string) string {
	emitter.GetState().AddValueToList(l, v)
	return ""
}

// RandomValueFromList returns a random value from Context list l
func RandomValueFromList(s string) any {
	return emitter.GetState().RandomValueFromList(s)
}

// RandomNValuesFromList returns a random value from Context list l
func RandomNValuesFromList(s string, n int) []string {
	l := emitter.GetState().RandomNValuesFromList(s, n)
	r := make([]string, 0)
	for i := range l {
		r = append(r, l[i].(string))
	}
	return r
}

// GetValuesFromList returns a value from Context list l at index
func GetValueFromListAtIndex(s string, index int) string {

	return emitter.GetState().GetValueFromListAtIndex(s, index).(string)
}

// GetV gets value s from Context
func GetV(s string) string {
	v, _ := emitter.GetState().Ctx.Load(s)
	return v.(string)
}

// SetV adds value v to Context
func SetV(s string, v string) string {
	emitter.GetState().Ctx.Store(s, v)
	return ""
}

// FromCsv gets the label value from csv file
func FromCSV(c string) string {
	return emitter.GetState().FromCSV(c)
}
