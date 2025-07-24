package middlewares

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/metric"
)

func RequestsCounter(counter metric.Int64Counter, duration metric.Float64Histogram) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := time.Now()
			counter.Add(r.Context(), 1)
			next.ServeHTTP(w, r)
			duration.Record(r.Context(), time.Since(n).Seconds())
		})
	}
}
