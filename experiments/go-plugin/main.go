package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-plugin"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
)

type MyProducer struct{}

func (MyProducer) Produce(key []byte, v []byte, headers map[string]string) error {
	value := []byte(fmt.Sprintf("%s\n\nWritten from plugin-go-grpc", string(v)))
	file, err := os.Create("kv_" + string(key))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(value)
	for k, v := range headers {
		file.Write([]byte(fmt.Sprintf("%s: %s\n", k, v)))
	}

	return nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: jrpc.Handshake,
		Plugins: map[string]plugin.Plugin{
			"producer": &jrpc.ProducerGRPCPlugin{Impl: &MyProducer{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
