package pkg

import "github.com/pkg/errors"

var (
	ErrInvalidSharesThreshold = errors.New("minimum shares threshold cannot be larger than share count")
	ErrNilSharedFile          = errors.New("nil shared file")
)
