package storage

import "github.com/pkg/errors"

var (
	ErrShareNotExists   = errors.New("share doesn't exist in store")
	ErrMoreThanOneMatch = errors.New("found more than one match for share")
)
