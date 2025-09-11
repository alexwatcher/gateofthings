package profiles

import (
	"context"

	profilesv1 "github.com/alexwatcher/gateofthings/protos/gen/go/profiles/v1"
	"google.golang.org/grpc"
)

type Profiles interface {
	Create(ctx context.Context, email string) (eml string, err error)
	GetMe(ctx context.Context, email string) (token string, err error)
}

type serverAPI struct {
	profilesv1.UnimplementedProfilesServer
	profiles Profiles
}

func Register(gRPC *grpc.Server, profiles Profiles) {
	profilesv1.RegisterProfilesServer(gRPC, &serverAPI{profiles: profiles})
}

func (s *serverAPI) Create(ctx context.Context, in *profilesv1.CreateRequest) (*profilesv1.CreateResponse, error) {

	return nil, nil
}

func (s *serverAPI) GetMe(ctx context.Context, in *profilesv1.GetMeRequest) (*profilesv1.GetMeResponse, error) {

	return nil, nil
}
