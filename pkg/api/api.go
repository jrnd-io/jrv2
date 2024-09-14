package api

import (
	"bytes"
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

func UserTemplateList() []*TemplateInfo {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.UserDir, "templates"))
	return templateList(templateDir)
}

func SystemTemplateList() []*TemplateInfo {
	templateDir := os.ExpandEnv(fmt.Sprintf("%s/%s", constants.SystemDir, "templates"))
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
			valid, err := isValidTemplate(t)
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

func isValidTemplate(t []byte) (bool, error) {

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
