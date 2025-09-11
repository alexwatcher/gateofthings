package profiles

import (
	profilesv1 "github.com/alexwatcher/gateofthings/protos/gen/go/profiles/v1"
	"google.golang.org/grpc"
)

type Profiles interface {
}

type serverAPI struct {
	profilesv1.UnimplementedProfilesServer
	profiles Profiles
}

func Register(gRPC *grpc.Server, auth Profiles) {
	profilesv1.RegisterProfilesServer(gRPC, &serverAPI{profiles: auth})
}
