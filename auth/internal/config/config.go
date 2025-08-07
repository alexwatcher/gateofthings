package config

import (
	"time"

	scfg "github.com/alexwatcher/gateofthings/shared/pkg/config"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env       string               `env:"ENV" envDefault:"local"`
	Secret    string               `env:"SECRET,required"`
	TokenTTL  time.Duration        `env:"TOKEN_TTL,required"`
	Telemetry scfg.TelemetryConfig `envPrefix:"TELEMETRY_"`
	GRPC      scfg.GRPCConfig      `envPrefix:"GRPC_"`
	Database  scfg.DatabaseConfig  `envPrefix:"DB_"`
}

// MustLoad loads configuration from environment variables into a Config instance.
// If the environment variables can't be parsed, it panics with the error.
func MustLoad() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	return env.Must(&cfg, err)
}
