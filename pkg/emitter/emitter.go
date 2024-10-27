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

package emitter

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ugol/uticker/t"

	"github.com/jrnd-io/jrv2/pkg/function"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	"github.com/jrnd-io/jrv2/pkg/tpl"
)

var Emitters map[string][]Config

type Throughput float64

type Ticker struct {
	Type           string
	ImmediateStart bool
	Num            int
	Frequency      time.Duration
	Duration       time.Duration
	Throughput     Throughput
	// Parameters     map[string]any `mapstructure:",remain"`
}

type Emitter struct {
	Config *Config

	StopChannel    chan struct{}
	Ticker         *t.UTicker
	KeyTemplate    *tpl.Tpl
	ValueTemplate  *tpl.Tpl
	HeaderTemplate *tpl.Tpl
	OutputTemplate *tpl.Tpl

	plugin *plugin.Plugin
}

func NewFromConfig(cfg Config) (*Emitter, error) {

	log.Debug().Msg("new emitter from config")
	e := &Emitter{
		Config:      &cfg,
		StopChannel: make(chan struct{}),
	}
	log.Debug().Interface("config", cfg).Msg("Creating new emitter from config")

	if err := e.SetTemplates(); err != nil {
		return nil, err
	}
	return e, nil
}
func New(options ...func(*Emitter)) (*Emitter, error) {

	tc := &Ticker{
		Type:           "simple",
		ImmediateStart: false,
		Num:            DefaultNum,
		Frequency:      DefaultFrequency,
		Duration:       DefaultDuration,
		Throughput:     DefaultThroughput,
	}

	e := &Emitter{
		Config: &Config{
			Tick:           *tc,
			Preload:        DefaultPreloadSize,
			Name:           DefaultEmitterName,
			Locale:         DefaultLocale,
			Output:         DefaultOutput,
			KeyTemplate:    DefaultKeyTemplate,
			ValueTemplate:  DefaultValueTemplate,
			HeaderTemplate: DefaultHeaderTemplate,
			OutputTemplate: DefaultOutputTemplate,
			Oneline:        false,
		},
		StopChannel: make(chan struct{}),
	}

	for _, option := range options {
		option(e)
	}

	return e, nil
}

func (e *Emitter) SetPlugin(plugin *plugin.Plugin) {
	e.plugin = plugin
}

func (e *Emitter) SetTemplates() error {

	var vTplText string
	var err error
	if e.Config.Embedded {
		vTplText = e.Config.ValueTemplate
	} else {
		vTplText, err = tpl.GetRawTemplate(e.Config.ValueTemplate)
		if err != nil {
			return err
		}
	}
	valueTpl, err := tpl.New("value", vTplText, function.Map())
	if err != nil {
		return err
	}
	e.ValueTemplate = valueTpl
	return nil

}

func (e *Emitter) Produce(ctx context.Context, key []byte, value []byte, headers map[string]string) (*jrpc.ProduceResponse, error) {
	if e.plugin == nil {
		return nil, errors.New("emitter plugin not initialized")
	}
	return e.plugin.Produce(ctx, key, value, headers)
}

func (e *Emitter) StartTicker() {
	e.Ticker = t.NewUTicker(t.WithFrequency(e.Config.Tick.Frequency))
}

func WithName(n string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.Name = n
	}
}

func WithLocale(l string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.Locale = l
	}
}

func WithOutput(o string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.Output = o
	}
}

func WithOneline(o bool) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.Oneline = o
	}
}

func WithImmediateStart(i bool) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.Tick.ImmediateStart = i
	}
}

func WithNum(n int) func(*Emitter) {
	if n < 1 {
		log.Warn().Msg("Num should be at least 1, setting to default")
		n = DefaultNum
	}
	return func(e *Emitter) {
		e.Config.Tick.Num = n
	}
}

func WithFrequency(f time.Duration) func(*Emitter) {
	if f <= 0 {
		log.Warn().Msg("Frequency is <=0, setting to default")
		f = DefaultFrequency
	}
	return func(e *Emitter) {
		e.Config.Tick.Frequency = f
	}
}

func WithThroughput(t string) func(*Emitter) {
	throughput, err := ParseThroughput(t)
	if err != nil {
		log.Error().Err(err).Msg("Error in parsing throughput, setting to default")
		throughput = DefaultThroughput
	}
	return func(e *Emitter) {
		e.Config.Tick.Throughput = throughput
	}
}

func WithDuration(d time.Duration) func(*Emitter) {
	if d <= 0 {
		log.Warn().Msg("Duration is <=0, setting to default")
		d = DefaultDuration
	}
	return func(e *Emitter) {
		e.Config.Tick.Duration = d
	}
}

func WithPreload(n int) func(*Emitter) {
	if n < 0 {
		log.Warn().Msg("Preload should be positive, setting to default")
		n = DefaultPreloadSize
	}
	return func(e *Emitter) {
		e.Config.Preload = n
	}
}

func WithKeyTemplate(k string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.KeyTemplate = k
	}
}

func WithValueTemplate(v string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.ValueTemplate = v
	}
}

func WithHeaderTemplate(h string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.HeaderTemplate = h
	}
}

func WithOutputTemplate(o string) func(*Emitter) {
	return func(e *Emitter) {
		e.Config.OutputTemplate = o
	}
}

func CalculateFrequency(bytes int, num int, throughput Throughput) time.Duration {

	if throughput == 0 {
		return 0
	}
	totalBytes := float64(bytes) * float64(num)

	// Calculate the frequency in milliseconds
	frequency := (totalBytes / float64(throughput)) * 1000

	return time.Duration(frequency) * time.Millisecond
}

//gocyclo:ignore
func ParseThroughput(input string) (Throughput, error) {

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
		return Throughput(value / 8), nil
	case "B":
		return Throughput(value), nil
	case "kb", "Kb":
		return Throughput(value * 1024 / 8), nil
	case "mb", "Mb":
		return Throughput(value * 1024 * 1024 / 8), nil
	case "gb", "Gb":
		return Throughput(value * 1024 * 1024 * 1024 / 8), nil
	case "tb", "Tb":
		return Throughput(value * 1024 * 1024 * 1024 * 1024 / 8), nil
	case "kB", "KB":
		return Throughput(value * 1024), nil
	case "mB", "MB":
		return Throughput(value * 1024 * 1024), nil
	case "gB", "GB":
		return Throughput(value * 1024 * 1024 * 1024), nil
	case "tB", "TB":
		return Throughput(value * 1024 * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unsupported unit: %s", unitStr)
	}
}
