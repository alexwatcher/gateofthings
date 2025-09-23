package postgresql

import (
	"context"

	"github.com/alexwatcher/gateofthings/profiles/internal/models"
	"github.com/jackc/pgx/v5"
)

type ProfilesRepo struct {
	conn *pgx.Conn
}

func NewProfilesRepo(conn *pgx.Conn) *ProfilesRepo {
	return &ProfilesRepo{conn: conn}
}

func (r *ProfilesRepo) Insert(ctx context.Context, profile *models.Profile) (string, error) {
	return "", nil
}

func (r *ProfilesRepo) Get(ctx context.Context, id string) (*models.Profile, error) {
	return nil, nil
}
