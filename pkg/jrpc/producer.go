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
	Produce([]byte, []byte, map[string]string) (*ProduceResponse, error)
}

// GRPC Client
type GRPCClient struct{ client ProducerClient }

func (m *GRPCClient) Produce(key []byte, v []byte, headers map[string]string) (*ProduceResponse, error) {
	resp, err := m.client.Produce(context.Background(), &ProduceRequest{Key: key, Value: v, Headers: headers})
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

	return m.Impl.Produce(req.Key, req.Value, req.Headers)
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
