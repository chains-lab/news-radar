package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/news-radar/internal/api"
	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/workers/eventlistener"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, domain service.Domain, cfg *config.Config) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	run(func() { api.Run(ctx, cfg, domain) })

	run(func() { eventlistener.Listen(ctx, cfg, domain) })
}
