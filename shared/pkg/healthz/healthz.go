package healthz

import (
	"context"
	"fmt"
	"net/http"
)

// Probe defines a function that performs a health check and returns an error if it fails.
type Probe func(context.Context) error

// Option configures a HealthCheck instance.
type Option interface {
	apply(*HealthCheck)
}

// HealthCheck represents a set of probes for Kubernetes health checks.
type HealthCheck struct {
	port  uint16
	ready Probe
	live  Probe
}

// New creates a new HealthCheck with the provided options.
func New(options ...Option) *HealthCheck {
	hc := &HealthCheck{
		port: 3000,
	}
	for _, o := range options {
		o.apply(hc)
	}
	return hc
}

// Run starts an HTTP server exposing health endpoints for Kubernetes probes.
func (c *HealthCheck) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	c.register(mux)

	server := &http.Server{Addr: fmt.Sprintf(":%d", c.port), Handler: mux}

	// Gracefully shutdown when the context is canceled.
	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

// MustRun starts an HTTP server exposing health endpoints.
// panic in case of error
func (c *HealthCheck) MustRun(ctx context.Context) {
	if err := c.Run(ctx); err != nil {
		panic(err)
	}
}

// RegisterHandlers registers all health endpoints on the given ServeMux.
// If nil is passed, handlers are registered on http.DefaultServeMux.
func (c *HealthCheck) register(mux *http.ServeMux) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	c.addHandler(mux, "/readyz", c.ready, "ready", "not ready")
	c.addHandler(mux, "/livez", c.live, "live", "not live")
}

// addHandler registers a single probe endpoint with appropriate HTTP status codes.
func (c *HealthCheck) addHandler(mux *http.ServeMux, path string, probe Probe, okMsg, failMsg string) {
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if probe != nil && probe(r.Context()) != nil {
			http.Error(w, failMsg, http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(okMsg))
	})
}

// withProbeOption is a generic option type for applying probes.
type withValueOption struct {
	setter func(*HealthCheck, any)
	value  any
}

func (o withValueOption) apply(hc *HealthCheck) {
	o.setter(hc, o.value)
}

// WithReadyProbe sets the readiness probe.
func WithReadyProbe(probe Probe) Option {
	return withValueOption{setter: func(h *HealthCheck, p any) { h.ready = p.(Probe) }, value: probe}
}

// WithLiveProbe sets the liveness probe.
func WithLiveProbe(probe Probe) Option {
	return withValueOption{setter: func(h *HealthCheck, p any) { h.live = p.(Probe) }, value: probe}
}

// WithPort sets port to run on http server.
func WithPort(port uint16) Option {
	return withValueOption{setter: func(h *HealthCheck, p any) { h.port = p.(uint16) }, value: port}
}
