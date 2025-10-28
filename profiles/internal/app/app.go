package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"

	grpcprofiles "github.com/alexwatcher/gateofthings/profiles/internal/grpc/profiles"
	"github.com/alexwatcher/gateofthings/profiles/internal/repository/postgresql"
	"github.com/alexwatcher/gateofthings/profiles/internal/services"
	"github.com/alexwatcher/gateofthings/shared/pkg/config"
	"github.com/alexwatcher/gateofthings/shared/pkg/grpc/interceptors/metadataextractor"
	"github.com/alexwatcher/gateofthings/shared/pkg/grpc/interceptors/tracing"
	"github.com/alexwatcher/gateofthings/shared/pkg/grpc/interceptors/valid"
	sharedpgsql "github.com/alexwatcher/gateofthings/shared/pkg/repository/postgresql"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	gRPConfig  config.GRPCSrvConfig
	dbConn     *pgx.Conn
	isRunning  int32
}

// New initializes a new instance of the App struct with a gRPC server
// listening on the specified port. It registers the authentication
// service with the server and returns the configured App instance.
func New(ctx context.Context, gRPConfig config.GRPCSrvConfig, dbConfig config.DatabaseConfig) *App {

	dbConn, err := sharedpgsql.NewConnection(ctx, dbConfig)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		panic(err)
	}
	repo := postgresql.NewProfilesRepo(dbConn)

	profilesService := services.NewProfiles(repo)
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			tracing.TracingInterceptor(),
			valid.UnaryInterceptor,
			metadataextractor.ExtractMetadataInterceptor,
		),
	)
	grpcprofiles.Register(gRPCServer, profilesService)
	return &App{
		gRPCServer: gRPCServer,
		gRPConfig:  gRPConfig,
		dbConn:     dbConn,
	}
}

// MustRun starts the gRPC server and panics if it can't be started.
func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

// Run starts the gRPC server and logs the port it is listening on. If the
// server can't be started, it returns an error.
func (a *App) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.gRPConfig.Port))
	if err != nil {
		return err
	}
	atomic.StoreInt32(&a.isRunning, 1)
	defer atomic.StoreInt32(&a.isRunning, 0)
	slog.Info("gRPC server started", "port", a.gRPConfig.Port)
	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("app.run: %w", err)
	}
	return nil
}

// Stop gracefully stops the gRPC server, ensuring that it no longer accepts new connections
// and waits for all ongoing RPCs to complete before shutting down.
func (a *App) Stop(ctx context.Context) {
	slog.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}

// Readiness probe
func (a *App) Ready() error {
	err := a.dbConn.Ping(context.Background())
	if err != nil {
		slog.Warn("app: live probe failed", "error", err)
		return err
	}
	if atomic.LoadInt32(&a.isRunning) == 0 {
		slog.Warn("app: live probe failed: app is not running")
		return fmt.Errorf("app: is not running")
	}
	return nil
}
