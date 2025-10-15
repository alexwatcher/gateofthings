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
	Upsert(ctx context.Context, id string, profile *models.Profile) (string, error)
	Get(ctx context.Context, id string) (*models.Profile, error)
}

type Profiles struct {
	repo ProfilesRepo
}

// NewProfiles initializes a new instance of the Profiles service with the given repository
func NewProfiles(repo ProfilesRepo) *Profiles {
	return &Profiles{repo: repo}
}

func (p *Profiles) UpdateMyProfile(ctx context.Context, profile *models.Profile) (string, error) {
	op := "profiles.updatemyprofile"
	log := slog.With("op", op)

	userId := contextutils.XUserIdFromContext(ctx)
	if userId == "" {
		log.Error("x-user-id not specified")
		return "", fmt.Errorf("%s: %w", op, models.ErrUnauthenticated)
	}

	log.Info("update profile")
	profileId, err := p.repo.Upsert(ctx, userId, profile)
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

func (p *Profiles) GetMyProfile(ctx context.Context) (*models.Profile, error) {
	op := "profiles.getmyprofile"
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
			log.Info("profile not found, return default")
			return &models.Profile{Id: userId, IsProvisioned: false}, nil
		}
		log.Error("failed get profile", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("profile found")
	profile.IsProvisioned = true
	return profile, nil
}
