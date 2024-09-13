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
	_ "embed"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

//go:embed descriptions.yaml
var descriptions []byte

func GetDescription(name string) (Description, bool) {
	v, ok := funcDesc[name]
	return v, ok
}

func DescriptionMap() map[string]Description {
	return funcDesc
}

var funcDesc map[string]Description

func init() {
	err := yaml.Unmarshal(descriptions, &funcDesc)
	if err != nil {
		log.Fatal().Err(err).Msg("error unmarshalling function descriptions")
	}

}

type Description struct {
	Name        string
	Category    string
	Description string
	Parameters  string
	Localizable string
	Return      string
	Example     string
	Output      string
}
