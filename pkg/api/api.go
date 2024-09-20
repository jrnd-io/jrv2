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

package api

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jrnd-io/jrv2/pkg/constants"
	"github.com/jrnd-io/jrv2/pkg/function"
)

func WithName(n string) func(*Emitter) {
	return func(e *Emitter) {
		e.Name = n
	}
}

func WithLocale(l string) func(*Emitter) {
	return func(e *Emitter) {
		e.Locale = l
	}
}

func WithImmediateStart(i bool) func(*Emitter) {
	return func(e *Emitter) {
		e.Tick.ImmediateStart = i
	}
}

func WithNum(n int) func(*Emitter) {
	if n < 1 {
		panic("JR should generate at least 1 object per iteration")
	}
	return func(e *Emitter) {
		e.Tick.Num = n
	}
}

func WithFrequency(f time.Duration) func(*Emitter) {
	if f <= 0 {
		panic("non-positive interval for Frequency")
	}
	return func(e *Emitter) {
		e.Tick.Frequency = f
	}
}

func WithDuration(d time.Duration) func(*Emitter) {
	if d <= 0 {
		panic("non-positive interval for NewTicker")
	}
	return func(e *Emitter) {
		e.Tick.Duration = d
	}
}

func WithPreload(n int) func(*Emitter) {
	if n < 0 {
		panic("Preload should be positive")
	}
	return func(e *Emitter) {
		e.Preload = n
	}
}

func WithKeyTemplate(k string) func(*Emitter) {
	return func(e *Emitter) {
		e.KeyTemplate = k
	}
}

func WithValueTemplate(v string) func(*Emitter) {
	return func(e *Emitter) {
		e.ValueTemplate = v
	}
}

func WithHeaderTemplate(h string) func(*Emitter) {
	return func(e *Emitter) {
		e.HeaderTemplate = h
	}
}

func WithOutputTemplate(o string) func(*Emitter) {
	return func(e *Emitter) {
		e.OutputTemplate = o
	}
}

func NewEmitter(options ...func(*Emitter)) (*Emitter, error) {

	defaultDuration, err := time.ParseDuration(constants.DefaultDuration)
	if err != nil {
		return nil, err
	}
	defaultFrequency, err := time.ParseDuration(constants.DefaultFrequency)
	if err != nil {
		return nil, err
	}

	t := &Ticker{
		Type:           "simple",
		ImmediateStart: false,
		Num:            constants.DefaultNum,
		Frequency:      defaultFrequency,
		Duration:       defaultDuration,
	}

	e := &Emitter{
		Tick:           *t,
		Preload:        constants.DefaultPreloadSize,
		Name:           constants.DefaultEmitterName,
		Locale:         constants.DefaultLocale,
		KeyTemplate:    constants.DefaultKeyTemplate,
		ValueTemplate:  constants.DefaultValueTemplate,
		HeaderTemplate: constants.DefaultHeaderTemplate,
		OutputTemplate: constants.DefaultOutputTemplate,
	}

	for _, option := range options {
		option(e)
	}

	return e, nil
}

func GetRawTemplate(name string) (string, error) {
	return getTemplate(name)
}

func GetRawValidatedTemplate(name string) (string, error) {

	t, err := getTemplate(name)
	if err != nil {
		return "", err
	}
	valid, err := IsValidTemplate(t)

	if !valid || err != nil {
		return "", errors.New("invalid template")
	}
	return t, nil
}

func IsValidTemplate(t string) (bool, error) {

	tt, err := template.New("test").Funcs(function.Map()).Parse(t)
	if err != nil {
		return false, err
	}

	var buf bytes.Buffer
	if err = tt.Execute(&buf, nil); err != nil {
		return false, err
	}

	return true, err

}

func SystemTemplateList() []*TemplateInfo {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.JR_SystemDir, "templates"))
	return templateList(templateDir)
}

func UserTemplateList() []*TemplateInfo {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.JR_UserDir, "templates"))
	return templateList(templateDir)
}

func getTemplate(name string) (string, error) {
	systemTemplateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.JR_SystemDir, "templates"))
	userTemplateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.JR_UserDir, "templates"))
	templateScript, err := os.ReadFile(fmt.Sprintf("%s/%s.tpl", userTemplateDir, name))
	if err != nil {
		templateScript, err = os.ReadFile(fmt.Sprintf("%s/%s.tpl", systemTemplateDir, name))
		if err != nil {
			return "", err
		}
	}
	return string(templateScript), nil
}

func templateList(templateDir string) []*TemplateInfo {

	templateList := make([]*TemplateInfo, 0)

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return templateList
	}

	_ = filepath.Walk(templateDir, func(path string, info os.FileInfo, _ error) error {
		if !info.IsDir() && strings.HasSuffix(path, "tpl") {

			t, _ := os.ReadFile(path)
			valid, err := IsValidTemplate(string(t))
			name, _ := strings.CutSuffix(info.Name(), ".tpl")
			templateInfo := TemplateInfo{
				Name:     name,
				IsValid:  valid,
				FullPath: path,
				Error:    err,
			}
			templateList = append(templateList, &templateInfo)
		}
		return nil
	})
	return templateList
}

/*
func countFilesInDir(dir string) int {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdirnames(-1)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
	return len(list)
}
*/
