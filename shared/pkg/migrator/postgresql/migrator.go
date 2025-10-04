package postgresql

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/alexwatcher/gateofthings/shared/pkg/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(cfg config.DatabaseConfig) {
	sslMode := "disable"
	if cfg.SSLMode {
		sslMode = "enable"
	}
	postgreDatasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, sslMode)

	db, err := sql.Open("pgx", postgreDatasource)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		panic(err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("failed to set dialect", "error", err)
		panic(err)
	}

	if err := goose.Up(db, cfg.Migrations); err != nil {
		slog.Error("failed to apply migrations", "error", err)
		panic(err)
	}

	slog.Info("migrations applied successfully.")
}
