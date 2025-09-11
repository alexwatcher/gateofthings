package services

import "context"

type Profiles struct {
	repo ProfilesRepo
}

type ProfilesRepo interface {
}

// NewProfiles initializes a new instance of the Profiles service with the given repository
func NewProfiles(repo ProfilesRepo) *Profiles {
	return &Profiles{repo: repo}
}

func (p *Profiles) Create(ctx context.Context, email string) (eml string, err error) {
	return "", nil
}

func (p *Profiles) GetMe(ctx context.Context, email string) (token string, err error) {
	return "", nil
}
