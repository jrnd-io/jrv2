package loop

import (
	"context"
	"fmt"
	"github.com/jrnd-io/jrv2/pkg/state"
	"os"
	"os/signal"
	"sync"

	"github.com/jrnd-io/jrv2/pkg/emitter"
	"github.com/rs/zerolog/log"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func DoLoop(ctx context.Context, emitters *orderedmap.OrderedMap[string, []emitter.Config]) {

	es := make([]*emitter.Emitter, 0)

	controlC, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

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
			es = append(es, emitter.NewFromConfig(cfg)) //nolint

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
							doTemplate(ctx, e)
						case <-e.StopChannel:
							return
						}

					}
				} else {
					log.Debug().
						Str("Emitter: %e", e.Config.Name).
						Msg("Exec do Template")
					doTemplate(ctx, e)
				}
			}(es[i])

		}
	}

	wg.Wait()
}

func doTemplate(ctx context.Context, em *emitter.Emitter) { //nolint
	// jrctx.JrContext.Locale = emitter.Locale
	// jrctx.JrContext.CountryIndex = functions.IndexOf(strings.ToUpper(emitter.Locale), "country")

	fmt.Printf("HERE: %s\n", em.Config.Name)
	for i := 0; i < em.Config.Tick.Num; i++ {
		state.GetState().Execution.CurrentIterationLoopIndex++

		// k := emitter.KTpl.Execute()
		// v := emitter.VTpl.Execute()
		if em.Config.Oneline { //nolint
			// v = strings.ReplaceAll(v, "\n", "")
		}
		// kInValue := functions.GetV("KEY")

		// if (kInValue) != "" {
		//	emitter.Producer.Produce(ctx, []byte(kInValue), []byte(v), nil)
		// } else {
		//	emitter.Producer.Produce(ctx, []byte(k), []byte(v), nil)
		// }

		// jrctx.JrContext.GeneratedObjects++
		// jrctx.JrContext.GeneratedBytes += int64(len(v))
	}

}
