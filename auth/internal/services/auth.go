package services

import (
	"context"
	"time"

	"github.com/alexwatcher/gateofthings/auth/internal/models"
)

type Auth struct {
	repo     UserRepo
	tokenTTL time.Duration
}

type UserRepo interface {
	InsertUser(ctx context.Context, email string, passhash []byte) error
	SelectUser(ctx context.Context, email string) (*models.User, error)
}

// NewAuth initializes a new instance of the Auth service with the given repository and
// token TTL.
func NewAuth(repo UserRepo, tokenTTL time.Duration) *Auth {
	return &Auth{repo: repo, tokenTTL: tokenTTL}
}

// Register creates a new user with the given email and password.
//
// If the user with the given email already exists, it returns an error.
func (s Auth) Register(ctx context.Context, email string, password string) (err error) {
	return nil
}

// Login generates a JWT token for the given user. If the user is not found or the
// password is incorrect, it returns an error.
func (s Auth) Login(ctx context.Context, login string, password string) (token string, err error) {
	return "", nil
}
