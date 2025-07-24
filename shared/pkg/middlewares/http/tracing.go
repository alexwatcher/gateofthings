package middlewares

import (
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

func Tracing(tracer trace.Tracer) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), "request")
			defer span.End()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
