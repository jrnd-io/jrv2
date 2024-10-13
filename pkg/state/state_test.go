package state_test

import (
	"github.com/jrnd-io/jrv2/pkg/random"
	"github.com/jrnd-io/jrv2/pkg/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetState(t *testing.T) {
	s := state.GetState()
	assert.NotNil(t, s)
	assert.Equal(t, "us", s.Locale)
}

func TestAddValueToList(t *testing.T) {
	s := state.GetState()
	s.AddValueToList("testList", "value1")
	s.AddValueToList("testList", "value2")

	value := s.RandomValueFromList("testList")
	assert.Contains(t, []string{"value1", "value2"}, value)
}

func TestCounter(t *testing.T) {
	s := state.GetState()

	// Test initial counter value
	assert.Equal(t, 10, s.Counter("testCounter", 10, 1))

	// Test increment
	assert.Equal(t, 11, s.Counter("testCounter", 10, 1))
	assert.Equal(t, 12, s.Counter("testCounter", 10, 1))
}

func TestRandomNValuesFromList(t *testing.T) {
	s := state.GetState()
	s.AddValueToList("testList", "value1")
	s.AddValueToList("testList", "value2")
	s.AddValueToList("testList", "value3")

	results := s.RandomNValuesFromList("testList", 2)
	assert.Len(t, results, 2)
	assert.Subset(t, []string{"value1", "value2", "value3"}, results)
}

func TestFromCSV(t *testing.T) {
	s := state.GetState()
	csvMap := state.CSVMap{
		0: {"column1": "value1", "column2": "value2"},
		1: {"column1": "value3", "column2": "value4"},
	}
	s.SetCSV(csvMap)

	assert.Equal(t, "value1", s.FromCSV("column1"))
	s.Execution.CurrentIterationLoopIndex++
	assert.Equal(t, "value3", s.FromCSV("column1"))
}

func TestValue(t *testing.T) {
	state := state.GetState()
	state.Ctx.Store("testKey", "testValue")

	assert.Equal(t, "testValue", state.Value("testKey"))
}

func init() {
	random.SetRandom(0)
}
