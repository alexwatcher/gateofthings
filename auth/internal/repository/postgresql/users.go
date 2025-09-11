package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/alexwatcher/gateofthings/auth/internal/models"
	"github.com/alexwatcher/gateofthings/auth/internal/repository"
	spgsql "github.com/alexwatcher/gateofthings/shared/pkg/repository/postgresql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersRepo struct {
	conn *pgx.Conn
}

func NewUsersRepo(conn *pgx.Conn) *UsersRepo {
	return &UsersRepo{conn: conn}
}

func (r *UsersRepo) Insert(ctx context.Context, email string, passhash []byte) (string, error) {
	op := "repository.users.Insert"

	id := uuid.New()

	sql, args, err := spgsql.SqlBuilder.
		Insert("auth_users").
		Columns("id", "email", "password_hash").
		Values(id, email, passhash).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == spgsql.ErrCodeUniqueViolation {
			return "", fmt.Errorf("%s: %w", op, repository.ErrUserAlreadyExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (r *UsersRepo) Get(ctx context.Context, email string) (models.User, error) {
	op := "repository.users.Select"

	var user models.User
	query, args, err := spgsql.SqlBuilder.
		Select("id", "email", "password_hash").
		From("auth_users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	err = r.conn.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
