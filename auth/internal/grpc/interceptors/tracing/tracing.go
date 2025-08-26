package tracing

import (
	"context"

	"github.com/alexwatcher/gateofthings/shared/pkg/telemetry/propagation"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TracingInterceptor() grpc.UnaryServerInterceptor {
	tracer := otel.Tracer("grpc-auth-tracer")
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.GRPCMetadataCarrier(md))
		}

		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		return handler(ctx, req)
	}
}
