package app

import (
	"fmt"
	"log/slog"
	"net"

	grpcauth "github.com/alexwatcher/gateofthings/auth/internal/grpc/auth"
	"github.com/alexwatcher/gateofthings/shared/pkg/config"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	gRPConfig  config.GRPCConfig
}

// New initializes a new instance of the App struct with a gRPC server
// listening on the specified port. It registers the authentication
// service with the server and returns the configured App instance.
func New(gRPConfig config.GRPCConfig) *App {
	gRPCServer := grpc.NewServer()
	grpcauth.Register(gRPCServer)
	return &App{
		gRPCServer: gRPCServer,
		gRPConfig:  gRPConfig,
	}
}

// MustRun starts the gRPC server and panics if it can't be started.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run starts the gRPC server and logs the port it is listening on. If the
// server can't be started, it returns an error.
func (a *App) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.gRPConfig.Port))
	if err != nil {
		return err
	}
	slog.Info("gRPC server started", "port", a.gRPConfig.Port)
	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("app.run: %w", err)
	}
	return nil
}

// Stop gracefully stops the gRPC server, ensuring that it no longer accepts new connections
// and waits for all ongoing RPCs to complete before shutting down.
func (a *App) Stop() {
	slog.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
