package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/alexwatcher/gateofthings/gateway/internal/interceptors"
	"github.com/alexwatcher/gateofthings/gateway/internal/middlewares"
	"github.com/alexwatcher/gateofthings/gateway/internal/openapi"
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	"github.com/alexwatcher/gateofthings/shared/pkg/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	httpConfig config.HTTPSrvConfig
	authConfig config.GRPCClnConfig
	server     *http.Server
	openAPI    string
}

func New(ctx context.Context, httpConfig config.HTTPSrvConfig, authConfig config.GRPCClnConfig, openAPI string) *App {
	return &App{
		httpConfig: httpConfig,
		authConfig: authConfig,
		openAPI:    openAPI,
	}
}

// MustRun starts the HTTP server and panics if it can't be started.
func (a *App) MustRun(ctx context.Context) {
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

// Run starts the HTTP server and logs the port it is listening on. If the
// server can't be started, it returns an error.
func (a *App) Run(ctx context.Context) error {
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(interceptors.MakeSetSignInCookie()),
		runtime.WithMiddlewares(
			middlewares.TracingMiddleware,
			middlewares.MakeCSRFMiddleware([]string{"/v1/auth/signin", "/v1/auth/singup"}),
		),
	)

	opts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(interceptors.MakeTracingClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err := authv1.RegisterAuthHandlerFromEndpoint(ctx, mux, a.authConfig.Address, opts)
	if err != nil {
		return fmt.Errorf("app.run: register http endpoint: %w", err)
	}

	err = openapi.RegisteraOpenAPIEndpoint(mux, a.openAPI)
	if err != nil {
		return fmt.Errorf("app.run: register openapi endpoint: %w", err)
	}

	slog.Info("HTTP server started", "port", a.httpConfig.Port)
	a.server = &http.Server{Addr: fmt.Sprintf(":%d", a.httpConfig.Port), Handler: mux}
	if err := a.server.ListenAndServe(); err != nil && err != context.Canceled {
		return fmt.Errorf("app.run: %w", err)
	}
	return nil
}

// Stop gracefully stops the HTTP server, ensuring that it no longer accepts new connections
// and waits for all ongoing connections to complete before shutting down.
func (a *App) Stop(ctx context.Context) {
	slog.Info("stopping HTTP server")
	a.server.Shutdown(ctx)
}
