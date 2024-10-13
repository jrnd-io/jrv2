package emitter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/emitter"

	//	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestThroughput(t *testing.T) {
	t.Run("Throughput", func(t *testing.T) {
		throughput, err := emitter.ParseThroughput("200KB/s")
		assert.NoError(t, err)
		frequency := emitter.CalculateFrequency(1024, 2, throughput)
		fmt.Println(frequency.Milliseconds())
		expected := 10 * time.Millisecond
		assert.Equal(t, expected, frequency)
	})
}

func TestParseThroughput(t *testing.T) {
	t.Run("ParseThroughput", func(t *testing.T) {
		throughput := "2MB/s"
		TwoMegaBytesPerSecond := emitter.Throughput(2 * 1024 * 1024)
		actual, err := emitter.ParseThroughput(throughput)
		assert.NoError(t, err)
		assert.Equal(t, TwoMegaBytesPerSecond, actual)
		throughput = "500Kb/s"
		FiveHundredKilobitPerSecond := emitter.Throughput(500 * 1024 / 8)
		actual, err = emitter.ParseThroughput(throughput)
		assert.NoError(t, err)
		assert.Equal(t, FiveHundredKilobitPerSecond, actual)
	})
}

func TestNewEmitter(t *testing.T) {
	t.Run("Default values", func(t *testing.T) {
		em, err := emitter.New()
		assert.NoError(t, err)
		assert.NotNil(t, em)

		// Check default values
		assert.Equal(t, emitter.DefaultEmitterName, em.Config.Name)
		assert.Equal(t, emitter.DefaultLocale, em.Config.Locale)
		assert.Equal(t, emitter.DefaultPreloadSize, em.Config.Preload)
		assert.Equal(t, emitter.DefaultKeyTemplate, em.Config.KeyTemplate)
		assert.Equal(t, emitter.DefaultValueTemplate, em.Config.ValueTemplate)
		assert.Equal(t, emitter.DefaultHeaderTemplate, em.Config.HeaderTemplate)
		assert.Equal(t, emitter.DefaultOutputTemplate, em.Config.OutputTemplate)

		// Check Ticker default values
		assert.Equal(t, "simple", em.Config.Tick.Type)
		assert.False(t, em.Config.Tick.ImmediateStart)
		assert.Equal(t, emitter.DefaultNum, em.Config.Tick.Num)

		defaultFrequency := config.DefaultFrequency
		assert.Equal(t, defaultFrequency, em.Config.Tick.Frequency)

		defaultDuration := config.DefaultDuration
		assert.Equal(t, defaultDuration, em.Config.Tick.Duration)
	})

	t.Run("Custom values", func(t *testing.T) {
		customName := "CustomEmitter"
		customLocale := "fr"
		customPreload := 10
		customKeyTemplate := "custom_key_{{.ID}}"
		customValueTemplate := "custom_value_{{.Name}}"
		customHeaderTemplate := "custom_header"
		customOutputTemplate := "custom_output"
		customNum := 5
		customFrequency := 2 * time.Second
		customDuration := 1 * time.Minute

		em, err := emitter.New(
			emitter.WithName(customName),
			emitter.WithLocale(customLocale),
			emitter.WithPreload(customPreload),
			emitter.WithKeyTemplate(customKeyTemplate),
			emitter.WithValueTemplate(customValueTemplate),
			emitter.WithHeaderTemplate(customHeaderTemplate),
			emitter.WithOutputTemplate(customOutputTemplate),
			emitter.WithNum(customNum),
			emitter.WithFrequency(customFrequency),
			emitter.WithDuration(customDuration),
			emitter.WithImmediateStart(true),
		)

		assert.NoError(t, err)
		assert.NotNil(t, em)

		// Check custom values
		assert.Equal(t, customName, em.Config.Name)
		assert.Equal(t, customLocale, em.Config.Locale)
		assert.Equal(t, customPreload, em.Config.Preload)
		assert.Equal(t, customKeyTemplate, em.Config.KeyTemplate)
		assert.Equal(t, customValueTemplate, em.Config.ValueTemplate)
		assert.Equal(t, customHeaderTemplate, em.Config.HeaderTemplate)
		assert.Equal(t, customOutputTemplate, em.Config.OutputTemplate)

		// Check Ticker custom values
		assert.True(t, em.Config.Tick.ImmediateStart)
		assert.Equal(t, customNum, em.Config.Tick.Num)
		assert.Equal(t, customFrequency, em.Config.Tick.Frequency)
		assert.Equal(t, customDuration, em.Config.Tick.Duration)
	})
}
