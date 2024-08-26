package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
)

func main() {

	jrPlugin := os.Getenv("JR_PLUGIN")
	fmt.Println("JR_PLUGIN: ", jrPlugin)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  jrpc.Handshake,
		Plugins:          jrpc.PluginMap,
		SyncStdout:       os.Stdout, // sync stdout from plugin
		SyncStderr:       os.Stderr, // sync stderr from plugin
		Stderr:           os.Stderr,
		Cmd:              exec.Command("sh", "-c", jrPlugin),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Managed:          false,
		Logger: hclog.New(&hclog.LoggerOptions{ // disable logging
			Name:  "plugin",
			Level: hclog.Trace,
		}),
	})
	// defer client.Kill()
	//defer plugin.CleanupClients()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		panic(err)
	}
	defer rpcClient.Close()

	// Request the plugin
	raw, err := rpcClient.Dispense(jrpc.JRProducerGRPCPlugin)
	if err != nil {
		panic(err)
	}

	producer := raw.(jrpc.Producer)
	for i := 0; i < 10; i++ {
		resp, err := producer.Produce([]byte(fmt.Sprintf("pippo_%d", i)), []byte("pippo value"), map[string]string{"k1": "v1"})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Response: %v\n", resp)
	}

}
