package plugin_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

const (
	PluginName = "test_plugin"
)

func init() {
	plugin.RegisterLocalPlugin(PluginName, &plugin.Plugin{
		Producer: &MockProducer{},
	})
}

// MockProducer is a mock implementation of the Producer interface
type MockProducer struct {
}

func (m *MockProducer) Produce(_ context.Context, _ []byte, value []byte, _ map[string]string) (*jrpc.ProduceResponse, error) {
	return &jrpc.ProduceResponse{Bytes: uint64(len(value)), Message: "success"}, nil
}

func (m *MockProducer) Close() error {
	return nil
}

func TestPlugin(t *testing.T) {

	// Test New function
	p, err := plugin.New(PluginName, hclog.Debug)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	// Test Produce method
	ctx := context.Background()
	key := []byte("test_key")
	value := []byte("test_value")
	headers := map[string]string{"header1": "value1"}

	expectedResponse := &jrpc.ProduceResponse{Bytes: uint64(len(value)), Message: "success"}

	response, err := p.Produce(ctx, key, value, headers)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)

	// Test Close method
	err = p.Close()
	assert.NoError(t, err)
}
