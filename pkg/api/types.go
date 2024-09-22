package api

import (
	"text/template"
	"time"
)

const (
	bitMultiplier = 8
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

type Ticker struct {
	Type           string
	ImmediateStart bool
	Num            int
	Frequency      time.Duration
	Duration       time.Duration
	Parameters     map[string]any `mapstructure:",remain"`
}

type Emitter struct {
	Tick           Ticker `mapstructure:",squash"`
	Preload        int
	Name           string
	Locale         string
	KeyTemplate    string
	ValueTemplate  string
	HeaderTemplate string
	OutputTemplate string
	Output         string
	Oneline        bool
}

type Throughput float64
