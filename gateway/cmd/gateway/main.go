package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexwatcher/gateofthings/gateway/internal/app"
	"github.com/alexwatcher/gateofthings/gateway/internal/config"
	"github.com/alexwatcher/gateofthings/gateway/internal/consts"
	"github.com/alexwatcher/gateofthings/shared/pkg/telemetry"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	res := telemetry.MustCreateResource(consts.ServiceName, consts.ServiceVersion, cfg.Env)
	telemetry.MustInitLogger(context.Background(), res, cfg.Telemetry.LogsEndpoint)
	telemetry.MustInitTracer(context.Background(), res, cfg.Telemetry.TraceEndpoint)
	telemetry.MustInitMeter(context.Background(), res, cfg.Telemetry.MetricsEndpoint)

	slog.Info("starting application")
	application := app.New(ctx, cfg.HTTP, cfg.Auth, cfg.OpenAPI)
	go application.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	slog.Info("stopping application", "signal", sig)
	application.Stop(ctx)
}
