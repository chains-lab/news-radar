package cli

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/recovery-flow/comtools/logkit"
	"github.com/recovery-flow/news-radar/internal/app"
	"github.com/recovery-flow/news-radar/internal/config"
)

func Run(args []string) bool {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logkit.SetupLogger(cfg.Server.Log.Level, cfg.Server.Log.Format)
	logger.Info("Starting server...")

	var (
		application = kingpin.New("news-radar", "")
		runCmd      = application.Command("run", "run command")
		serviceCmd  = runCmd.Command("service", "run service")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	svc, err := app.NewApp(*cfg)
	if err != nil {
		logger.WithError(err).Error("failed to create svc")
		return false
	}

	cmd, err := application.Parse(args[1:])
	if err != nil {
		logger.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		runServices(ctx, &wg, svc, cfg)
	default:
		logger.Errorf("unknown command %s", cmd)
		return false
	}
	if err != nil {
		logger.WithError(err).Error("failed to exec cmd")
		return false
	}

	wgch := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgch)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Interrupt signal received: %v", ctx.Err())
		<-wgch
	case <-wgch:
		log.Print("All services stopped")
	}

	return true
}
