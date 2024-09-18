package api

import "time"

type Templates struct {
	TemplateList []*TemplateInfo
}

type TemplateInfo struct {
	Name     string
	IsValid  bool
	FullPath string
	Error    error
}

type Emitter struct {
	Name           string `mapstructure:"Name"`
	Locale         string `mapstructure:"Locale"`
	ImmediateStart bool
	Num            int           `mapstructure:"Num"`
	Frequency      time.Duration `mapstructure:"Frequency"`
	Duration       time.Duration `mapstructure:"Duration"`
	Preload        int           `mapstructure:"Preload"`
	KeyTemplate    string        `mapstructure:"KeyTemplate"`
	ValueTemplate  string        `mapstructure:"ValueTemplate"`
	HeaderTemplate string        `mapstructure:"HeaderTemplate"`
	OutputTemplate string        `mapstructure:"OutputTemplate"`
}
