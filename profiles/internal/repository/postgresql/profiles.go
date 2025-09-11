package postgresql

import "github.com/jackc/pgx/v5"

type ProfilesRepo struct {
	conn *pgx.Conn
}

func NewProfilesRepo(conn *pgx.Conn) *ProfilesRepo {
	return &ProfilesRepo{conn: conn}
}
