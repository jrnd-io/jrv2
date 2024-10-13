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

package plugin_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

type testProducer struct {
	plugin.Producer
	bytesToWrite uint64
	message      string
}

func (t *testProducer) Produce(_ context.Context, _ []byte, _ []byte, _ map[string]string) (*jrpc.ProduceResponse, error) {
	return &jrpc.ProduceResponse{
		Bytes:   t.bytesToWrite,
		Message: t.message,
	}, nil
}

func TestLocalPlugin(t *testing.T) {

	// create and register a local plugin
	pluginName := "test"
	bytesToWrite := uint64(100)
	message := "test message"
	p := &testProducer{
		bytesToWrite: bytesToWrite,
		message:      message,
	}
	plugin.RegisterLocalPlugin(pluginName, &plugin.Plugin{
		Name:     pluginName,
		Producer: p,
		IsRemote: false,
	})

	// get a client for the plugin
	c, err := plugin.New(pluginName, "", hclog.Off)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	resp, err := c.Produce(context.Background(), []byte("key"), []byte("value"), nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, bytesToWrite, resp.Bytes)
	assert.Equal(t, message, resp.Message)

}
