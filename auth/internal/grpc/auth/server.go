package auth

import (
	"context"
	"log/slog"

	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Register(context.Context, *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	slog.Warn("Register not implemented")
	return &authv1.RegisterResponse{}, nil
}

func (s *serverAPI) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	slog.Warn("Login not implemented")
	return &authv1.LoginResponse{
		Token: in.GetEmail() + in.GetPassword(),
	}, nil
}
