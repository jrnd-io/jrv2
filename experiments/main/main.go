package main

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
)

func main() {

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  jrpc.Handshake,
		Plugins:          jrpc.PluginMap,
		Cmd:              exec.Command("sh", "-c", os.Getenv("JR_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		panic(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(jrpc.JRProducerGRPCPlugin)
	if err != nil {
		panic(err)
	}

	producer := raw.(jrpc.Producer)
	producer.Produce([]byte("pippo"), []byte("pippo value"), map[string]string{"k1": "v1"})
	producer.Produce([]byte("pluto"), []byte("pluto value"), map[string]string{"k2": "v2"})

}
