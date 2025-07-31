package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/alexwatcher/gateofthings/auth/internal/lib/jwt"
	"github.com/alexwatcher/gateofthings/auth/internal/models"
	"github.com/alexwatcher/gateofthings/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	repo     UserRepo
	tokenTTL time.Duration
}

type UserRepo interface {
	Insert(ctx context.Context, email string, passhash []byte) (string, error)
	Get(ctx context.Context, email string) (models.User, error)
}

// NewAuth initializes a new instance of the Auth service with the given repository and
// token TTL.
func NewAuth(repo UserRepo, tokenTTL time.Duration) *Auth {
	return &Auth{repo: repo, tokenTTL: tokenTTL}
}

// Register creates a new user with the given email and password.
//
// If the user with the given email already exists, it returns an error.
func (s Auth) Register(ctx context.Context, email string, password string) (string, error) {
	op := "auth.register"
	log := slog.With("op", op, "email", email)

	log.Info("registering user")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.repo.Insert(ctx, email, hash)
	if err != nil {
		log.Error("failed to insert user", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user registered")

	return id, nil
}

// Login generates a JWT token for the given user. If the user is not found or the
// password is incorrect, it returns an error.
func (s Auth) Login(ctx context.Context, login string, password string) (string, error) {
	op := "auth.login"
	log := slog.With("op", op, "login", login)

	log.Info("logging in user")
	user, err := s.repo.Get(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Error("user not found", "error", err)
			return "", repository.ErrUserNotFound
		}

		log.Error("failed to select user", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		log.Error("failed to compare password hash", "error", err)
		return "", models.ErrInvalidCredentials
	}

	token, err := jwt.NewToken(user, s.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in")

	return token, nil
}
