package config

import "time"

type GRPCConfig struct {
	Port    int           `env:"PORT,required"`
	Timeout time.Duration `env:"TIMEOUT,required"`
}

type TelemetryConfig struct {
	TraceEndpoint   string `env:"TRACE"`
	MetricsEndpoint string `env:"METRICS"`
	LogsEndpoint    string `env:"LOGS"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	SSLMode  bool   `env:"SSL" envDefault:"false"`
	Name     string `env:"NAME,required"`
	User     string `env:"USER,required,unset"`
	Password string `env:"PASSWORD,required,unset"`
}
