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
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jrnd-io/jrv2/pkg/config"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/rs/zerolog/log"
)

const (
	EnvJRPlugin = "JR_PLUGIN"

	PythonPlugin       = "python"
	PythonPluginScript = "plugin.py"
)

type Plugin struct {
	Name      string
	RPCClient plugin.ClientProtocol
	Producer  Producer
	Command   string
	IsRemote  bool
}

func New(jrPlugin string, logLevel hclog.Level) (*Plugin, error) {

	var pl *Plugin
	var command string

	if jrPlugin != PythonPlugin {
		pl = pluginMap[jrPlugin]
		if pl == nil {
			log.Error().Str("plugin", jrPlugin).Msg("plugin not found")
			return nil, fmt.Errorf("plugin %s not found", jrPlugin)
		}
		log.Debug().
			Str("plugin", jrPlugin).
			Interface("pluginmap", pluginMap).
			Bool("remote", pl.IsRemote).
			Msg("initialize plugin")

		if !pl.IsRemote {
			return &Plugin{
				Name:     jrPlugin,
				Producer: pl.Producer,
			}, nil
		}
		command = pl.Command
	} else {
		command = "python"
	}

	if logLevel == 0 {
		logLevel = hclog.Off
	}

	cmd := "sh"
	args := []string{
		"-c",
	}

	pCmd := sanitize(command)
	pArgs := ""

	configFile := getConfigFileFor(jrPlugin)
	if configFile != "" {
		log.Debug().Str("configFile", configFile).Msg("adding configuration file to args")
		pArgs = fmt.Sprintf("--config '%s'", configFile)
	}

	if jrPlugin == PythonPlugin {
		pArgs = PythonPluginScript
	}
	pCmd = fmt.Sprintf("%s %s", pCmd, pArgs)

	args = append(args, pCmd)

	log.Debug().
		Str("cmd", cmd).
		Strs("args", args).
		Str("plugin", jrPlugin).
		Msg("creating client")
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: jrpc.Handshake,
		Plugins:         jrpc.PluginMap,
		SyncStdout:      os.Stdout, // sync stdout from plugin
		SyncStderr:      os.Stderr, // sync stderr from plugin
		Stderr:          os.Stderr,
		//TODO: make this portable
		Cmd:              exec.Command(cmd, args...), //nolint
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Managed:          false,
		Logger: hclog.New(&hclog.LoggerOptions{
			Name:  "jrplugin",
			Level: logLevel,
		}),
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(jrpc.JRProducerGRPCPlugin)
	if err != nil {
		return nil, err
	}

	p, ok := raw.(jrpc.Producer)
	if !ok {
		return nil, err
	}

	return &Plugin{
		Name:     jrPlugin,
		Producer: NewAdapter(p),
		IsRemote: true,
		Command:  command,
	}, nil
}

func (c *Plugin) Produce(ctx context.Context, key []byte, value []byte, headers map[string]string, configParams map[string]string) (*jrpc.ProduceResponse, error) {

	return c.Producer.Produce(ctx, key, value, headers, configParams)

}

func (c *Plugin) Close() error {
	if c.RPCClient != nil {
		return c.RPCClient.Close()
	}
	return nil
}

func sanitize(c string) string {
	return strings.ReplaceAll(c, " ", "\\ ")
}

func getConfigFileFor(jrPlugin string) string {
	systemConfig := fmt.Sprintf("%s%cplugins%c%s.conf.json", config.JrSystemDir, os.PathSeparator, os.PathSeparator, jrPlugin)
	userConfig := fmt.Sprintf("%s%cplugins%c%s.conf.json", config.JrUserDir, os.PathSeparator, os.PathSeparator, jrPlugin)
	log.Debug().
		Str("systemConfig", systemConfig).
		Str("userConfig", userConfig).
		Str("plugin", jrPlugin).
		Msg("searching config")
	if _, err := os.Stat(userConfig); err == nil {
		log.Debug().Str("plugin", jrPlugin).Str("config", userConfig).Msg("found config")
		return userConfig
	}
	if _, err := os.Stat(systemConfig); err == nil {
		log.Debug().Str("plugin", jrPlugin).Str("config", systemConfig).Msg("found config")
		return systemConfig
	}
	log.Debug().Str("plugin", jrPlugin).Msg("config not found, proceeding without")
	return ""
}
