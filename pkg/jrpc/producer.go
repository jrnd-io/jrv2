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

type Producer interface {
	Produce([]byte, []byte, map[string]string) error
}
type GRPCClient struct{ client ProducerClient }

func (m *GRPCClient) Produce(key []byte, v []byte, headers map[string]string) error {
	_, err := m.client.Produce(context.Background(), &ProduceRequest{Key: key, Value: v, Headers: headers})
	return err
}

type GRPCServer struct {
	Impl Producer
}

func (m *GRPCServer) Produce(
	_ context.Context,
	req *ProduceRequest) (*Empty, error) {

	return &Empty{}, m.Impl.Produce(req.Key, req.Value, req.Headers)
}

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
