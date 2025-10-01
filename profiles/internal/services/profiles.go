package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/alexwatcher/gateofthings/profiles/internal/models"
	"github.com/alexwatcher/gateofthings/shared/pkg/contextutils"
)

type ProfilesRepo interface {
	Insert(ctx context.Context, profile *models.Profile) (string, error)
	Get(ctx context.Context, id string) (*models.Profile, error)
}

type Profiles struct {
	repo ProfilesRepo
}

// NewProfiles initializes a new instance of the Profiles service with the given repository
func NewProfiles(repo ProfilesRepo) *Profiles {
	return &Profiles{repo: repo}
}

func (p *Profiles) Create(ctx context.Context, profile *models.Profile) (string, error) {
	op := "profiles.create"
	log := slog.With("op", op, "id", profile.Id, "name", profile.Name)

	userId := contextutils.XUserIdFromContext(ctx)
	if userId != profile.Id {
		log.Error("x-user-id not equalt to requested user id")
		return "", fmt.Errorf("%s: %w", op, models.ErrUnauthenticated)
	}

	log.Info("create profile")
	profileId, err := p.repo.Insert(ctx, profile)
	if err != nil {
		if errors.Is(err, models.ErrProfileAlreadyExists) {
			log.Error("profile already exists", "error", err)
			return "", fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed create profile", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("profile created")
	return profileId, nil
}

func (p *Profiles) GetMe(ctx context.Context) (*models.Profile, error) {
	op := "profiles.getme"
	log := slog.With("op", op)

	userId := contextutils.XUserIdFromContext(ctx)
	if userId == "" {
		log.Error("x-user-id not specified")
		return nil, fmt.Errorf("%s: %w", op, models.ErrUnauthenticated)
	}

	log.Info("get my profile")
	profile, err := p.repo.Get(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrProfileNotFound) {
			log.Error("profile not found", "error", err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("failed get profile", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("profile found")
	return profile, nil
}
