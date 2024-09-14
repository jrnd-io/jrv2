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
	"github.com/jrnd-io/jrv2/pkg/constants"
	"github.com/jrnd-io/jrv2/pkg/function"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateInfo struct {
	Name     string
	IsValid  bool
	FullPath string
	Error    error
}

func GetRawTemplate(name string) (string, error) {
	systemTemplateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.SystemDir, "templates"))
	userTemplateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.UserDir, "templates"))
	templateScript, err := os.ReadFile(fmt.Sprintf("%s/%s.tpl", userTemplateDir, name))
	if err != nil {
		templateScript, err = os.ReadFile(fmt.Sprintf("%s/%s.tpl", systemTemplateDir, name))
		if err != nil {
			return "", err
		}
	}
	valid, err := IsValidTemplate(templateScript)

	if !valid || err != nil {
		return "", errors.New("invalid template")
	} else {
		return string(templateScript), nil
	}
}

func IsValidTemplate(t []byte) (bool, error) {

	tt, err := template.New("test").Funcs(function.Map()).Parse(string(t))
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
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.SystemDir, "templates"))
	return templateList(templateDir)
}

func UserTemplateList() []*TemplateInfo {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.UserDir, "templates"))
	return templateList(templateDir)
}

func templateList(templateDir string) []*TemplateInfo {

	templateList := make([]*TemplateInfo, 0)

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return templateList
	}

	_ = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, "tpl") {

			t, _ := os.ReadFile(path)
			valid, err := IsValidTemplate(t)
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

func countFilesInDir(dir string) int {
	f, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	list, err := f.Readdirnames(-1)
	err = f.Close()

	if err != nil {
		panic(err)
	}
	return len(list)
}
