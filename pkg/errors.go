package pkg

import "github.com/pkg/errors"

var (
	ErrShareNotExists         = errors.New("share doesn't exist in store")
	ErrMoreThanOneMatch       = errors.New("found more than one match for share")
	ErrInvalidSharesThreshold = errors.New("minimum shares threshold cannot be larger than share count")
	ErrNilSharedFile          = errors.New("nil shared file")
)
