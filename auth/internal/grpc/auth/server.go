package auth

import (
	"context"
	"errors"

	"github.com/alexwatcher/gateofthings/auth/internal/models"
	"github.com/alexwatcher/gateofthings/auth/internal/repository"
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	SignUp(ctx context.Context, email string, password string) (id string, err error)
	SignIn(ctx context.Context, email string, password string) (token string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) SignUp(ctx context.Context, req *authv1.SignUpRequest) (*authv1.SignUpResponse, error) {
	id, err := s.auth.SignUp(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.SignUpResponse{Id: id}, nil
}

func (s *serverAPI) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	token, err := s.auth.SignIn(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authv1.SignInResponse{Token: token}, nil
}
