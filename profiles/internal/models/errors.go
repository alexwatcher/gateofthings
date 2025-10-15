package models

import "errors"

var (
	ErrProfileNotFound               = errors.New("profile not found")
	ErrProfileAlreadyExists          = errors.New("profile already exists")
	ErrProfilePropertiesNotSpecified = errors.New("profile properties not specified")
	ErrUnauthenticated               = errors.New("unauthenticated")
)
