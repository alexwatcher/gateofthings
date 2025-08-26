package middlewares

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel"
)

func TracingMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	tracer := otel.Tracer("grpc-gateway-tracer")
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		ctx, span := tracer.Start(r.Context(), r.URL.Path)
		defer span.End()

		next(w, r.WithContext(ctx), pathParams)
	}
}
