package metadataextractor

import (
	"context"

	"github.com/alexwatcher/gateofthings/shared/pkg/contextutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ExtractMetadataInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userid, ok := md["x-user-id"]; ok && len(userid) > 0 {
			return handler(contextutils.ContextWithXUserId(ctx, userid[0]), req)
		}
	}
	return handler(ctx, req)
}
