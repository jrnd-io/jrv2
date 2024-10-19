package loop

import (
	"context"
	"github.com/jrnd-io/jrv2/pkg/jrpc"
	"github.com/jrnd-io/jrv2/pkg/state"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/jrnd-io/jrv2/pkg/plugin"
	"github.com/rs/zerolog/log"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func DoLoop(ctx context.Context,
	pluginName string,
	pluginConfigFile string,
	emitters *orderedmap.OrderedMap[string, []emitter.Config]) error {

	// emitter slice
	es := make([]*emitter.Emitter, 0)

	//  ctrl-c signal
	controlC, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// wait group to synchronize tickers end
	var wg sync.WaitGroup

	// creating plugin
	p, err := plugin.New(pluginName, pluginConfigFile, hclog.Debug)
	defer func() {
		_ = p.Close()
	}()
	if err != nil {
		return err
	}

	// starting loop
	for e := emitters.Oldest(); e != nil; e = e.Next() {
		log.Debug().
			Str("emitter", e.Key).
			Int("len", len(e.Value)).
			Msg("Starting loop for emitters")

		for i, cfg := range e.Value {

			log.Debug().
				Int("emitter", i).
				Interface("config", cfg).
				Msg("Running emitter")
			em, err := emitter.NewFromConfig(cfg)
			if err != nil {
				return err
			}
			es = append(es, em) //nolint

			wg.Add(1)
			go func(e *emitter.Emitter) {
				defer wg.Done()

				frequency := e.Config.Tick.Frequency
				if frequency > 0 {
					log.Debug().
						Dur("frequency", frequency).
						Str("emitter", e.Config.Name).
						Msg("Starting ticker")
					//					ticker := time.NewTicker(frequency)
					//					defer ticker.Stop()
					e.StartTicker()

					for {
						select {
						case <-controlC.Done():
							stop()
							return
						case <-e.Ticker.C:
							doTemplate(ctx, p, e)
						case <-e.StopChannel:
							return
						}

					}
				} else {
					log.Debug().
						Str("Emitter: %e", e.Config.Name).
						Msg("Exec do Template")
					doTemplate(ctx, p, e)
				}
			}(es[i])

		}
	}

	wg.Wait()
	return nil
}

func doTemplate(ctx context.Context,
	p *plugin.Plugin,
	em *emitter.Emitter) { //nolint

	var err error

	localState := state.NewState()
	for i := 0; i < em.Config.Tick.Num; i++ {
		state.GetSharedState().Execution.CurrentIterationLoopIndex++

		keyText := ""
		valueText := ""

		if em.ValueTemplate != nil {
			valueText = em.ValueTemplate.ExecuteWith(localState)
			if em.Config.Oneline {
				valueText = strings.ReplaceAll(valueText, "\n", "")
			}
		}

		if em.KeyTemplate != nil {
			keyText = em.KeyTemplate.Execute()
			log.Debug().Str("key", keyText).Msg("key generated with template")
		} else {
			keyText = localState.Key
			log.Debug().Str("key", keyText).Msg("key generated within localState")
		}

		log.Debug().Str("plugin", p.Name).Msg("producing message")
		var resp *jrpc.ProduceResponse
		resp, err = p.Produce(ctx, []byte(keyText), []byte(valueText), localState.Header)
		if err != nil {
			log.Warn().
				Err(err).
				Str("plugin", p.Name).
				Msg("error in emission")
		}

		state.GetSharedState().Execution.GeneratedObjects++
		state.GetSharedState().Execution.GeneratedBytes += resp.Bytes
	}

}
