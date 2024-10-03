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
	"github.com/jrnd-io/jrv2/pkg/config"
	"github.com/jrnd-io/jrv2/pkg/loop"
	"github.com/jrnd-io/jrv2/pkg/types"
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

	emitters := orderedmap.New[string, []*types.Emitter](len(args))
	for _, name := range args {
		e := config.Emitters[name]
		if dryrun {
			fmt.Println("should set output to stdout")
		}
		emitters.Set(name, e)
	}
	RunEmitters(cmd.Context(), emitters)

}

func RunEmitters(ctx context.Context, emitters *orderedmap.OrderedMap[string, []*types.Emitter]) {
	// defer emitter.WriteStats()
	// defer emitter.CloseProducers(ctx, ems)
	// emittersToRun := emitter.Initialize(ctx, emitterNames, ems, dryrun)
	// emitter.DoLoop(ctx, emittersToRun)

	log.Debug().Msg("Running main loop")
	loop.DoLoop(ctx, emitters)

}

func init() {
	RunCmd.Flags().BoolP("dryrun", "d", false, "dryrun: output of the emitters to stdout")
}
