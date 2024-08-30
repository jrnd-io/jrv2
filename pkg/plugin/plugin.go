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
package plugin

import (
	"context"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/rs/zerolog/log"
)

const (
	EnvJRPlugin = "JR_PLUGIN"
)

type Plugin struct {
	rpcClient plugin.ClientProtocol
	producer  Producer
}

func New(jrPlugin string, logLevel hclog.Level) (*Plugin, error) {

	if localPluginMap[jrPlugin] != nil {
		return &Plugin{
			producer: localPluginMap[jrPlugin],
		}, nil
	}

	if logLevel == 0 {
		logLevel = hclog.Off
	}
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  jrpc.Handshake,
		Plugins:          jrpc.PluginMap,
		SyncStdout:       os.Stdout, // sync stdout from plugin
		SyncStderr:       os.Stderr, // sync stderr from plugin
		Stderr:           os.Stderr,
		Cmd:              exec.Command("sh", "-c", jrPlugin), //TODO: make this portable
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Managed:          false,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:  "jrplugin",
			Level: logLevel,
		}),
	})

	rpcClient, err := client.Client()
	if err != nil {
		log.Error().Err(err).Msg("failed to create client")
		return nil, err
	}

	raw, err := rpcClient.Dispense(jrpc.JRProducerGRPCPlugin)
	if err != nil {
		log.Error().Err(err).Msg("failed to dispense plugin")
		return nil, err
	}

	p, ok := raw.(jrpc.Producer)
	if !ok {
		log.Error().Msg("failed to cast plugin")
		return nil, err
	}

	return &Plugin{
		producer: NewAdapter(p),
	}, nil
}

func (c *Plugin) Produce(ctx context.Context, key []byte, value []byte, headers map[string]string) (*jrpc.ProduceResponse, error) {

	return c.producer.Produce(ctx, key, value, headers)

}

func (c *Plugin) Close() error {
	if c.rpcClient != nil {
		return c.rpcClient.Close()
	}
	return nil
}
