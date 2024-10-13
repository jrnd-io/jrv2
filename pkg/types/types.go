package types

import (
	"text/template"
)

type Templates struct {
	TemplateList []*TemplateInfo
}

type TemplateInfo struct {
	Name     string
	IsValid  bool
	FullPath string
	Error    error
	Template *template.Template
}
