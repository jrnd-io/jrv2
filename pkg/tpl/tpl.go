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

package tpl

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/jrnd-io/jrv2/pkg/types"
	"github.com/jrnd-io/jrv2/pkg/utils"
	"github.com/rs/zerolog/log"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Tpl struct {
	Template *template.Template
}

func New(name string, t string, fmap map[string]interface{}) (*Tpl, error) {

	tp, err := template.New(name).Funcs(fmap).Parse(t)
	if err != nil {
		return nil, err
	}

	tpl := &Tpl{
		Template: tp,
	}
	return tpl, nil
}

func (t *Tpl) Execute() string {
	return t.ExecuteWith(nil)
}

func (t *Tpl) ExecuteWith(data any) string {
	var buffer bytes.Buffer
	err := t.Template.Execute(&buffer, data)
	if err != nil {
		log.Fatal().Err(err).Msg("Error executing template")
	}
	return buffer.String()
}

func GetRawTemplate(name string) (string, error) {
	return getTemplate(name)
}

func GetRawValidatedTemplate(name string) (string, error) {

	t, err := getTemplate(name)
	if err != nil {
		return "", err
	}
	valid, _, err := IsValidTemplate(t)

	if !valid || err != nil {
		return "", errors.New("invalid template")
	}
	return t, nil
}

func IsValidTemplate(t string) (bool, *template.Template, error) {

	tt, err := template.New("test").Funcs(function.Map()).Parse(t)
	if err != nil {
		return false, nil, err
	}

	var buf bytes.Buffer
	if err = tt.Execute(&buf, nil); err != nil {
		return false, nil, err
	}

	return true, tt, err

}

func ExecuteTemplate(stringTemplate string, data any) (string, error) {
	tt, err := New("test", stringTemplate, function.Map())
	if err != nil {
		return "", err
	}
	return tt.ExecuteWith(data), nil
}

func ExecuteTemplateByName(name string, ctx any) (string, error) {
	t, err := getTemplate(name)
	if err != nil {
		return "", err
	}
	return ExecuteTemplate(t, ctx)
}

func SystemTemplateList() *orderedmap.OrderedMap[string, *types.TemplateInfo] {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", config.JrSystemDir, "templates"))
	return templateList(templateDir)
}

func UserTemplateList() *orderedmap.OrderedMap[string, *types.TemplateInfo] {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", config.JrUserDir, "templates"))
	return templateList(templateDir)
}

func getTemplate(name string) (string, error) {
	systemTemplateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", config.JrSystemDir, "templates"))
	userTemplateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", config.JrUserDir, "templates"))
	templateScript, err := os.ReadFile(fmt.Sprintf("%s/%s.tpl", userTemplateDir, name))
	if err != nil {
		templateScript, err = os.ReadFile(fmt.Sprintf("%s/%s.tpl", systemTemplateDir, name))
		if err != nil {
			return "", err
		}
	}
	return string(templateScript), nil
}

func templateList(templateDir string) *orderedmap.OrderedMap[string, *types.TemplateInfo] {

	howManyTemplatesInTemplateDir, _ := utils.CountFilesInDir(templateDir)
	templateList := orderedmap.New[string, *types.TemplateInfo](howManyTemplatesInTemplateDir)

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return templateList
	}

	_ = filepath.WalkDir(templateDir, func(path string, f fs.DirEntry, _ error) error {
		if strings.HasSuffix(path, "tpl") {

			t, _ := os.ReadFile(path)
			valid, tt, err := IsValidTemplate(string(t))
			name, _ := strings.CutSuffix(f.Name(), ".tpl")
			templateInfo := types.TemplateInfo{
				Name:     name,
				IsValid:  valid,
				FullPath: path,
				Template: tt,
				Error:    err,
			}
			templateList.Set(name, &templateInfo)
		}
		return nil
	})
	return templateList
}
