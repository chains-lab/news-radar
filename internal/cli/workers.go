package cli

import (
	"context"
	"sync"

	"github.com/recovery-flow/news-radar/internal/service"
	"github.com/recovery-flow/news-radar/internal/service/api"
	"github.com/recovery-flow/news-radar/internal/service/events/listener"
)

func runServices(ctx context.Context, wg *sync.WaitGroup, svc *service.Service) {
	run := func(f func()) {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	run(func() { listener.Listen(ctx, svc) })

	run(func() { api.Run(ctx, svc) })
}
