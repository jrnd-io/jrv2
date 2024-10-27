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

package jrpc

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	JRProducerGRPCPlugin = "jr_producer_grpc"
	JRMagicCookieKey     = "JR_PRODUCER_PLUGIN"
	JRMagicCookieValue   = "producer"
	JRProtocolVersion    = 1
)

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  JRProtocolVersion,
	MagicCookieKey:   JRMagicCookieKey,
	MagicCookieValue: JRMagicCookieValue,
}

// Interface
type Producer interface {
	Produce([]byte, []byte, map[string]string, map[string]string) (*ProduceResponse, error)
}

// GRPC Client
type GRPCClient struct{ client ProducerClient }

func (m *GRPCClient) Produce(key []byte, v []byte, headers map[string]string, configParams map[string]string) (*ProduceResponse, error) {
	resp, err := m.client.Produce(context.Background(), &ProduceRequest{Key: key, Value: v, Headers: headers, ConfigParams: configParams})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type GRPCServer struct {
	Impl Producer
}

func (m *GRPCServer) Produce(
	_ context.Context,
	req *ProduceRequest) (*ProduceResponse, error) {

	return m.Impl.Produce(req.Key, req.Value, req.Headers, req.ConfigParams)
}

// Plugin Map
var PluginMap = map[string]plugin.Plugin{
	JRProducerGRPCPlugin: &ProducerGRPCPlugin{},
}

type ProducerGRPCPlugin struct {
	plugin.Plugin
	Impl Producer
}

func (p *ProducerGRPCPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	RegisterProducerServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ProducerGRPCPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: NewProducerClient(c)}, nil
}
