package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/alexwatcher/gateofthings/profiles/internal/models"
	spgsql "github.com/alexwatcher/gateofthings/shared/pkg/repository/postgresql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProfilesRepo struct {
	conn *pgx.Conn
}

func NewProfilesRepo(conn *pgx.Conn) *ProfilesRepo {
	return &ProfilesRepo{conn: conn}
}

func (r *ProfilesRepo) Insert(ctx context.Context, profile *models.Profile) (string, error) {
	op := "repository.profiles.Insert"
	sql, args, err := spgsql.SqlBuilder.
		Insert("user_profiles").
		Columns("id", "name", "avatar").
		Values(profile.Id, profile.Name, profile.Avatar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == spgsql.ErrCodeUniqueViolation {
			return "", fmt.Errorf("%s: %w", op, models.ErrProfileAlreadyExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return profile.Id, nil
}

func (r *ProfilesRepo) Get(ctx context.Context, id string) (*models.Profile, error) {
	op := "repository.profiles.Get"

	var profile models.Profile
	query, args, err := spgsql.SqlBuilder.
		Select("id", "name", "avatar").
		From("user_profiles").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return &profile, fmt.Errorf("%s: %w", op, err)
	}

	err = r.conn.QueryRow(ctx, query, args...).Scan(&profile.Id, &profile.Name, &profile.Avatar)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &profile, fmt.Errorf("%s: %w", op, models.ErrProfileNotFound)
		}
		return &profile, fmt.Errorf("%s: %w", op, err)
	}

	return &profile, nil
}
