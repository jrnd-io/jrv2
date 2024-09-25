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
	"github.com/jrnd-io/jrv2/pkg/utils"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/jrnd-io/jrv2/pkg/tpl"
	"github.com/jrnd-io/jrv2/pkg/types"
	"github.com/rs/zerolog/log"

	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/function"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func WithName(n string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.Name = n
	}
}

func WithLocale(l string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.Locale = l
	}
}

func WithOutput(o string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.Output = o
	}
}

func WithOneline(o bool) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.Oneline = o
	}
}

func WithImmediateStart(i bool) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.Tick.ImmediateStart = i
	}
}

func WithNum(n int) func(*types.Emitter) {
	if n < 1 {
		log.Warn().Msg("Num should be at least 1, setting to default")
		n = config.DefaultNum
	}
	return func(e *types.Emitter) {
		e.Tick.Num = n
	}
}

func WithFrequency(f time.Duration) func(*types.Emitter) {
	if f <= 0 {
		log.Warn().Msg("Frequency is <=0, setting to default")
		f = config.DefaultFrequency
	}
	return func(e *types.Emitter) {
		e.Tick.Frequency = f
	}
}

func WithThroughput(t string) func(*types.Emitter) {
	throughput, err := ParseThroughput(t)
	if err != nil {
		log.Error().Err(err).Msg("Error in parsing throughput, setting to default")
		throughput = config.DefaultThroughput
	}
	return func(e *types.Emitter) {
		e.Tick.Throughput = throughput
	}
}

func WithDuration(d time.Duration) func(*types.Emitter) {
	if d <= 0 {
		log.Warn().Msg("Duration is <=0, setting to default")
		d = config.DefaultDuration
	}
	return func(e *types.Emitter) {
		e.Tick.Duration = d
	}
}

func WithPreload(n int) func(*types.Emitter) {
	if n < 0 {
		log.Warn().Msg("Preload should be positive, setting to default")
		n = config.DefaultPreloadSize
	}
	return func(e *types.Emitter) {
		e.Preload = n
	}
}

func WithKeyTemplate(k string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.KeyTemplate = k
	}
}

func WithValueTemplate(v string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.ValueTemplate = v
	}
}

func WithHeaderTemplate(h string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.HeaderTemplate = h
	}
}

func WithOutputTemplate(o string) func(*types.Emitter) {
	return func(e *types.Emitter) {
		e.OutputTemplate = o
	}
}

func NewEmitter(options ...func(*types.Emitter)) (*types.Emitter, error) {

	t := &types.Ticker{
		Type:           "simple",
		ImmediateStart: false,
		Num:            config.DefaultNum,
		Frequency:      config.DefaultFrequency,
		Duration:       config.DefaultDuration,
		Throughput:     config.DefaultThroughput,
	}

	e := &types.Emitter{
		Tick:           *t,
		Preload:        config.DefaultPreloadSize,
		Name:           config.DefaultEmitterName,
		Locale:         config.DefaultLocale,
		Output:         config.DefaultOutput,
		KeyTemplate:    config.DefaultKeyTemplate,
		ValueTemplate:  config.DefaultValueTemplate,
		HeaderTemplate: config.DefaultHeaderTemplate,
		OutputTemplate: config.DefaultOutputTemplate,
		Oneline:        false,
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

func ExecuteTemplate(stringTemplate string, ctx any) (string, error) {
	tt, err := tpl.New("test", stringTemplate, function.Map(), ctx)
	if err != nil {
		return "", err
	}
	return tt.Execute(), nil
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

func CalculateFrequency(bytes int, num int, throughput types.Throughput) time.Duration {

	if throughput == 0 {
		return 0
	}
	totalBytes := float64(bytes) * float64(num)

	// Calculate the frequency in milliseconds
	frequency := (totalBytes / float64(throughput)) * 1000

	return time.Duration(frequency) * time.Millisecond
}

//gocyclo:ignore
func ParseThroughput(input string) (types.Throughput, error) {

	if input == "" {
		return -1, nil
	}
	re := regexp.MustCompile(`^((?:0|[1-9]\d*)(?:\.?\d*))([KkMmGgTt][Bb])/([smhd])$`)
	match := re.FindStringSubmatch(input)

	if len(match) != 4 {
		return 0, fmt.Errorf("invalid input format: %s", input)
	}

	valueStr := match[1]
	unitStr := match[2]
	timeStr := match[3]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numeric value: %w", err)
	}

	switch timeStr {
	case "s":
		// nothing to do
	case "m":
		value /= 60
	case "h":
		value /= 3600
	case "d":
		value /= 86400
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unitStr)
	}

	switch unitStr {
	case "b":
		return types.Throughput(value / 8), nil
	case "B":
		return types.Throughput(value), nil
	case "kb", "Kb":
		return types.Throughput(value * 1024 / 8), nil
	case "mb", "Mb":
		return types.Throughput(value * 1024 * 1024 / 8), nil
	case "gb", "Gb":
		return types.Throughput(value * 1024 * 1024 * 1024 / 8), nil
	case "tb", "Tb":
		return types.Throughput(value * 1024 * 1024 * 1024 * 1024 / 8), nil
	case "kB", "KB":
		return types.Throughput(value * 1024), nil
	case "mB", "MB":
		return types.Throughput(value * 1024 * 1024), nil
	case "gB", "GB":
		return types.Throughput(value * 1024 * 1024 * 1024), nil
	case "tB", "TB":
		return types.Throughput(value * 1024 * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unitStr)
	}
}
