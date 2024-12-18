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
	"bytes"
	"github.com/jrnd-io/jrv2/pkg/state"
	"regexp"
	"text/template"

	"github.com/rs/zerolog/log"
)

func ExecuteTemplate(key *template.Template, value *template.Template, oneline bool) (string, string, error) {

	var kBuffer, vBuffer bytes.Buffer
	var err error

	if err = key.Execute(&kBuffer, state.GetSharedState()); err != nil {
		log.Error().Err(err).Msg("Error executing key template")
	}
	k := kBuffer.String()

	if err = value.Execute(&vBuffer, state.GetSharedState()); err != nil {
		log.Error().Err(err).Msg("Error executing value template")
	}
	v := vBuffer.String()

	if oneline {
		re := regexp.MustCompile(`\r?\n?`)
		v = re.ReplaceAllString(v, "")
	}

	// TODO: maybe this does not go here since the actual produced bytes are in the plugin response
	state.GetSharedState().Execution.GeneratedObjects++
	state.GetSharedState().Execution.GeneratedBytes += uint64(len(v))

	return k, v, err
}
