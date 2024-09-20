package emitter_test

import (
	"testing"

	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/stretchr/testify/assert"
)

func TestGetState(t *testing.T) {
	state := emitter.GetState()
	assert.NotNil(t, state)
	assert.Equal(t, "us", state.Locale)
}

func TestAddValueToList(t *testing.T) {
	state := emitter.GetState()
	state.AddValueToList("testList", "value1")
	state.AddValueToList("testList", "value2")

	value := state.RandomValueFromList("testList")
	assert.Contains(t, []string{"value1", "value2"}, value)
}

func TestCounter(t *testing.T) {
	state := emitter.GetState()

	// Test initial counter value
	assert.Equal(t, 10, state.Counter("testCounter", 10, 1))

	// Test increment
	assert.Equal(t, 11, state.Counter("testCounter", 10, 1))
	assert.Equal(t, 12, state.Counter("testCounter", 10, 1))
}

func TestRandomNValuesFromList(t *testing.T) {
	state := emitter.GetState()
	state.AddValueToList("testList", "value1")
	state.AddValueToList("testList", "value2")
	state.AddValueToList("testList", "value3")

	results := state.RandomNValuesFromList("testList", 2)
	assert.Len(t, results, 2)
	assert.Subset(t, []string{"value1", "value2", "value3"}, results)
}

func TestFromCSV(t *testing.T) {
	state := emitter.GetState()
	csvMap := emitter.CSVMap{
		0: {"column1": "value1", "column2": "value2"},
		1: {"column1": "value3", "column2": "value4"},
	}
	state.SetCSV(csvMap)

	assert.Equal(t, "value1", state.FromCSV("column1"))
	state.Execution.CurrentIterationLoopIndex++
	assert.Equal(t, "value3", state.FromCSV("column1"))
}

func TestValue(t *testing.T) {
	state := emitter.GetState()
	state.Ctx.Store("testKey", "testValue")

	assert.Equal(t, "testValue", state.Value("testKey"))
}
