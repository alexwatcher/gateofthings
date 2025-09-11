package services

type Profiles struct {
	repo ProfilesRepo
}

type ProfilesRepo interface {
}

// NewProfiles initializes a new instance of the Profiles service with the given repository
func NewProfiles(repo ProfilesRepo) *Profiles {
	return &Profiles{repo: repo}
}
