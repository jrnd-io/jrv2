package loop

import (
	"context"
	"github.com/jrnd-io/jrv2/pkg/types"
	"github.com/ugol/uticker/t"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"os"
	"os/signal"
	"sync"
	"time"
)

func DoLoop(ctx context.Context, emitters *orderedmap.OrderedMap[string, []*types.Emitter]) {

	numEmitters := emitters.Len()
	es := make([]*types.Emitter, numEmitters)
	timers := make([]*t.UTicker, numEmitters)
	stopChannels := make([]chan struct{}, numEmitters)

	for e := emitters.Oldest(); e != nil; e = e.Next() {
		for i, v := range e.Value {
			es[i] = v
			stopChannels[i] = make(chan struct{})
			timers[i] = t.NewUTicker(t.WithFrequency(v.Tick.Duration))
		}
	}

	var wg sync.WaitGroup
	wg.Add(numEmitters)

	controlC, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i := 0; i < numEmitters; i++ {

		index := i

		go func(timerIndex int) {
			defer wg.Done()

			frequency := es[timerIndex].Tick.Frequency
			if frequency > 0 {
				ticker := time.NewTicker(frequency)
				defer ticker.Stop()
				for {
					select {
					case <-controlC.Done():
						stop()
						return
					case <-ticker.C:
						doTemplate(ctx, es[index])
					case <-stopChannels[timerIndex]:
						return
					}

				}
			} else {
				doTemplate(ctx, es[index])
			}
		}(index)

	}

	wg.Wait()
}

func doTemplate(ctx context.Context, emitter *types.Emitter) { //nolint
	// jrctx.JrContext.Locale = emitter.Locale
	// jrctx.JrContext.CountryIndex = functions.IndexOf(strings.ToUpper(emitter.Locale), "country")

	for i := 0; i < emitter.Tick.Num; i++ {
		// jrctx.JrContext.CurrentIterationLoopIndex++

		// k := emitter.KTpl.Execute()
		// v := emitter.VTpl.Execute()
		if emitter.Oneline { //nolint
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
