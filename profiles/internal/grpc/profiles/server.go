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
	UpdateMyProfile(ctx context.Context, properties *models.Profile) (id string, err error)
	GetMyProfile(ctx context.Context) (profile *models.Profile, err error)
}

type serverAPI struct {
	profilesv1.UnimplementedProfilesServer
	profiles Profiles
}

func Register(gRPC *grpc.Server, profiles Profiles) {
	profilesv1.RegisterProfilesServer(gRPC, &serverAPI{profiles: profiles})
}

func (s *serverAPI) UpdateMyProfile(ctx context.Context, in *profilesv1.UpdateMyProfileRequest) (*profilesv1.UpdateMyProfileResponse, error) {
	if in.Properties == nil {
		return nil, status.Error(codes.Internal, models.ErrProfilePropertiesNotSpecified.Error())
	}

	properties := &models.Profile{
		Name:   in.Properties.Name,
		Avatar: in.Properties.Avatar,
	}
	_, err := s.profiles.UpdateMyProfile(ctx, properties)
	if err != nil {
		if errors.Is(err, models.ErrProfileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &profilesv1.UpdateMyProfileResponse{}, nil
}

func (s *serverAPI) GetMyProfile(ctx context.Context, in *profilesv1.GetMyProfileRequest) (*profilesv1.GetMyProfileResponse, error) {
	profile, err := s.profiles.GetMyProfile(ctx)
	if err != nil {
		if errors.Is(err, models.ErrUnauthenticated) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &profilesv1.GetMyProfileResponse{Porfile: &profilesv1.Profile{
		Id:            profile.Id,
		IsProvisioned: profile.IsProvisioned,
		Properties: &profilesv1.ProfileProperties{
			Name:   profile.Name,
			Avatar: profile.Avatar,
		},
	}}, nil
}
