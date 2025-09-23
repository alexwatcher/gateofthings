package profiles

import (
	"context"
	"errors"

	"github.com/alexwatcher/gateofthings/profiles/internal/models"
	profilesv1 "github.com/alexwatcher/gateofthings/protos/gen/go/profiles/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Profiles interface {
	Create(ctx context.Context, profile *models.Profile) (id string, err error)
	GetMe(ctx context.Context) (profile *models.Profile, err error)
}

type serverAPI struct {
	profilesv1.UnimplementedProfilesServer
	profiles Profiles
}

func Register(gRPC *grpc.Server, profiles Profiles) {
	profilesv1.RegisterProfilesServer(gRPC, &serverAPI{profiles: profiles})
}

func (s *serverAPI) Create(ctx context.Context, in *profilesv1.CreateRequest) (*profilesv1.CreateResponse, error) {
	_, err := s.profiles.Create(ctx, &models.Profile{})
	if err != nil {
		if errors.Is(err, models.ErrProfileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &profilesv1.CreateResponse{}, nil
}

func (s *serverAPI) GetMe(ctx context.Context, in *profilesv1.GetMeRequest) (*profilesv1.GetMeResponse, error) {
	_, err := s.profiles.GetMe(ctx)
	if err != nil {
		if errors.Is(err, models.ErrProfileNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &profilesv1.GetMeResponse{}, nil
}
