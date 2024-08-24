package main

import (
	"context"
	"jrclient/proto"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "producer",
}

type Producer interface {
	Produce(key []byte, v []byte) error
}
type GRPCClient struct{ client proto.ProducerClient }

func (m *GRPCClient) Produce(key []byte, v []byte) error {
	_, err := m.client.Produce(context.Background(), &proto.ProduceRequest{Key: key, Value: v})
	return err
}

type GRPCServer struct {
	// This is the real implementation
	Impl Producer
}

func (m *GRPCServer) Produce(
	ctx context.Context,
	req *proto.ProduceRequest) (*proto.Empty, error) {

	return &proto.Empty{}, m.Impl.Produce(req.Key, req.Value)
}

var PluginMap = map[string]plugin.Plugin{
	"producer_grpc": &ProducerGRPCPlugin{},
}

type ProducerGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl Producer
}

func (p *ProducerGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProducerServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *ProducerGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewProducerClient(c)}, nil
}
