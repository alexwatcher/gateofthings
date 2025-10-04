package models

import "errors"

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrUnauthenticated      = errors.New("unauthenticated")
)
