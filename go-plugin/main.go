package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
)

type MyProducer struct{}

func (MyProducer) Produce(key []byte, v []byte) error {
	value := []byte(fmt.Sprintf("%s\n\nWritten from plugin-go-grpc", string(v)))
	return os.WriteFile("kv_"+string(key), value, 0644)
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"producer": &ProducerGRPCPlugin{Impl: &MyProducer{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
