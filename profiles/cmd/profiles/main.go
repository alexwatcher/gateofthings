package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexwatcher/gateofthings/profiles/internal/app"
	"github.com/alexwatcher/gateofthings/profiles/internal/config"
	"github.com/alexwatcher/gateofthings/profiles/internal/consts"
	sharedpgsql "github.com/alexwatcher/gateofthings/shared/pkg/migrator/postgresql"
	"github.com/alexwatcher/gateofthings/shared/pkg/telemetry"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	res := telemetry.MustCreateResource(consts.ServiceName, consts.ServiceVersion, cfg.Env)
	telemetry.MustInitLogger(context.Background(), res, cfg.Telemetry.LogsEndpoint)
	telemetry.MustInitTracer(context.Background(), res, cfg.Telemetry.TraceEndpoint)
	telemetry.MustInitMeter(context.Background(), res, cfg.Telemetry.MetricsEndpoint)

	slog.Info("start migration")
	sharedpgsql.Migrate(cfg.Database)
	slog.Info("end migration")

	slog.Info("starting application")
	application := app.New(ctx, cfg.GRPC, cfg.Database)
	go application.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	slog.Info("stopping application", "signal", sig)
	application.Stop(ctx)
}
