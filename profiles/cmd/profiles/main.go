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
	"github.com/alexwatcher/gateofthings/shared/pkg/healthz"
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

	application := app.New(ctx, cfg.GRPC, cfg.Database)

	hc := healthz.New(
		healthz.WithPort(cfg.HealthPort),
		healthz.WithLiveProbe(func(ctx context.Context) error { return nil }),
		healthz.WithReadyProbe(func(ctx context.Context) error { return application.Ready() }),
	)
	go hc.MustRun(ctx)

	slog.Info("start migration")
	sharedpgsql.Migrate(cfg.Database)
	slog.Info("end migration")

	slog.Info("starting application")
	go application.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	slog.Info("stopping application", "signal", sig)
	application.Stop(ctx)
}
