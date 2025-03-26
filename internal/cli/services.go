package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/news-radar/internal/api"
	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/services/eventlistener"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, app app.App, cfg *config.Config) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	API := api.NewAPI(cfg, app)
	run(func() { API.Run(ctx) })

	run(func() { eventlistener.Listen(ctx, cfg, app) })
}
