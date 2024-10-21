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

package emitter

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/loop"
	"github.com/jrnd-io/jrv2/pkg/plugin/local/console"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all or selected configured emitters",
	Long:  `Run all or selected configured emitters`,
	Args:  cobra.MinimumNArgs(1),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	dryrun, _ := cmd.Flags().GetBool("dryrun")
	pluginName, err := cmd.Flags().GetString("output")
	if err != nil {
		log.Warn().Err(err).Msgf("error in getting producer flag, defaulting to %s", console.PluginName)
		pluginName = console.PluginName
	}

	var logLevel hclog.Level
	verbosity, err := cmd.PersistentFlags().GetCount("output-log-level")
	if err != nil {
		verbosity = 0
	}

	switch verbosity {
	case 1:
		logLevel = hclog.Debug
	case 2:
		logLevel = hclog.Info
	case 3:
		logLevel = hclog.Warn
	case 4:
		logLevel = hclog.Error
	default:
		logLevel = hclog.Off
	}

	emitters := orderedmap.New[string, []emitter.Config](len(args))
	for _, name := range args {
		e := emitter.Emitters[name]
		if dryrun {
			pluginName = console.PluginName
			fmt.Println("should set output to stdout")
		}
		emitters.Set(name, e)
	}
	RunEmitters(cmd.Context(), pluginName, emitters, logLevel)

}

func RunEmitters(ctx context.Context,
	pluginName string,
	emitters *orderedmap.OrderedMap[string, []emitter.Config],
	pluginLogLevel hclog.Level) {

	log.Debug().Msg("Running main loop")
	if err := loop.DoLoop(ctx, emitters, pluginName, pluginLogLevel); err != nil {
		fmt.Printf("%v\n", err)
	}

}

func init() {
	RunCmd.Flags().BoolP("dryrun", "d", false, "dryrun: output of the emitters to stdout")
	RunCmd.Flags().StringP("output", "o", "", "name of output producer")
	RunCmd.Flags().CountP("output-log-level", "l", "name of output producer")
}
