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
	secret   string
}

type UserRepo interface {
	Insert(ctx context.Context, email string, passhash []byte) (string, error)
	Get(ctx context.Context, email string) (models.User, error)
}

// NewAuth initializes a new instance of the Auth service with the given repository and
// token TTL.
func NewAuth(repo UserRepo, secret string, tokenTTL time.Duration) *Auth {
	return &Auth{repo: repo, secret: secret, tokenTTL: tokenTTL}
}

// SignUp creates a new user with the given email and password.
//
// If the user with the given email already exists, it returns an error.
func (s Auth) SignUp(ctx context.Context, email string, password string) (string, error) {
	op := "auth.signup"
	log := slog.With("op", op, "email", email)

	log.Info("sign up user")
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
	log.Info("user signed up")

	return id, nil
}

// SignIn generates a JWT token for the given user. If the user is not found or the
// password is incorrect, it returns an error.
func (s Auth) SignIn(ctx context.Context, email string, password string) (string, error) {
	op := "auth.signin"
	log := slog.With("op", op, "email", email)

	log.Info("sign in user")
	user, err := s.repo.Get(ctx, email)
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

	token, err := jwt.NewToken(user, s.secret, s.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user signed in")

	return token, nil
}
