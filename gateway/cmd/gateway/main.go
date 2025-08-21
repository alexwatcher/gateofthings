package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexwatcher/gateofthings/gateway/internal/app"
	"github.com/alexwatcher/gateofthings/gateway/internal/config"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	application := app.New(ctx, cfg.HTTP, cfg.Auth)
	go application.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	slog.Info("stopping application", "signal", sig)
	application.Stop(ctx)
}
