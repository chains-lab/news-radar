package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/news-radar/internal/config"
	"github.com/recovery-flow/news-radar/internal/service/api"
	"github.com/recovery-flow/news-radar/internal/service/app"
	"github.com/recovery-flow/news-radar/internal/service/workers/evelisten"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, app app.App, cfg *config.Config) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	run(func() { api.Run(ctx, cfg, app) })

	run(func() { evelisten.Listen(ctx, cfg, app) })
}
