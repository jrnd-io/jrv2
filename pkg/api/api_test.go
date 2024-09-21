package api_test

import (
	"testing"
	"time"

	"github.com/jrnd-io/jrv2/pkg/api"
	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewEmitter(t *testing.T) {
	t.Run("Default values", func(t *testing.T) {
		emitter, err := api.NewEmitter()
		assert.NoError(t, err)
		assert.NotNil(t, emitter)

		// Check default values
		assert.Equal(t, config.DefaultEmitterName, emitter.Name)
		assert.Equal(t, config.DefaultLocale, emitter.Locale)
		assert.Equal(t, config.DefaultPreloadSize, emitter.Preload)
		assert.Equal(t, config.DefaultKeyTemplate, emitter.KeyTemplate)
		assert.Equal(t, config.DefaultValueTemplate, emitter.ValueTemplate)
		assert.Equal(t, config.DefaultHeaderTemplate, emitter.HeaderTemplate)
		assert.Equal(t, config.DefaultOutputTemplate, emitter.OutputTemplate)

		// Check Ticker default values
		assert.Equal(t, "simple", emitter.Tick.Type)
		assert.False(t, emitter.Tick.ImmediateStart)
		assert.Equal(t, config.DefaultNum, emitter.Tick.Num)

		defaultFrequency, _ := time.ParseDuration(config.DefaultFrequency)
		assert.Equal(t, defaultFrequency, emitter.Tick.Frequency)

		defaultDuration, _ := time.ParseDuration(config.DefaultDuration)
		assert.Equal(t, defaultDuration, emitter.Tick.Duration)
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

		emitter, err := api.NewEmitter(
			api.WithName(customName),
			api.WithLocale(customLocale),
			api.WithPreload(customPreload),
			api.WithKeyTemplate(customKeyTemplate),
			api.WithValueTemplate(customValueTemplate),
			api.WithHeaderTemplate(customHeaderTemplate),
			api.WithOutputTemplate(customOutputTemplate),
			api.WithNum(customNum),
			api.WithFrequency(customFrequency),
			api.WithDuration(customDuration),
			api.WithImmediateStart(true),
		)

		assert.NoError(t, err)
		assert.NotNil(t, emitter)

		// Check custom values
		assert.Equal(t, customName, emitter.Name)
		assert.Equal(t, customLocale, emitter.Locale)
		assert.Equal(t, customPreload, emitter.Preload)
		assert.Equal(t, customKeyTemplate, emitter.KeyTemplate)
		assert.Equal(t, customValueTemplate, emitter.ValueTemplate)
		assert.Equal(t, customHeaderTemplate, emitter.HeaderTemplate)
		assert.Equal(t, customOutputTemplate, emitter.OutputTemplate)

		// Check Ticker custom values
		assert.True(t, emitter.Tick.ImmediateStart)
		assert.Equal(t, customNum, emitter.Tick.Num)
		assert.Equal(t, customFrequency, emitter.Tick.Frequency)
		assert.Equal(t, customDuration, emitter.Tick.Duration)
	})
}