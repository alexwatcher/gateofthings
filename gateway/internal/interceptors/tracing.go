package interceptors

import (
	"context"

	"github.com/alexwatcher/gateofthings/shared/pkg/telemetry/propagation"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func MakeTracingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			md = md.Copy()
		} else {
			md = metadata.New(nil)
		}
		otel.GetTextMapPropagator().Inject(ctx, propagation.GRPCMetadataCarrier(md))
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
