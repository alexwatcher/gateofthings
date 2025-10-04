package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/alexwatcher/gateofthings/shared/pkg/config"
	"github.com/jackc/pgx/v5"
)

var SqlBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewConnection(ctx context.Context, cfg config.DatabaseConfig) (*pgx.Conn, error) {
	sslMode := "disable"
	if cfg.SSLMode {
		sslMode = "enable"
	}
	postgreDatasource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, sslMode)

	conn, err := pgx.Connect(ctx, postgreDatasource)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
