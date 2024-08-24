package main

import (
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

func main() {

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  Handshake,
		Plugins:          PluginMap,
		Cmd:              exec.Command("sh", "-c", os.Getenv("KV_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		panic(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("producer_grpc")
	if err != nil {
		panic(err)
	}

	producer := raw.(Producer)
	producer.Produce([]byte("pippo"), []byte("pippo value"))
	producer.Produce([]byte("pluto"), []byte("pluto value"))

}
