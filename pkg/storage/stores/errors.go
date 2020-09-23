package stores

import "github.com/pkg/errors"

var (
	ErrShareNotExists   = errors.New("share doesn't exist in store")
	ErrMoreThanOneMatch = errors.New("found more than one match for share")

	// Missing/invalid config attributes errors.
	ErrInvalidStoreType         = errors.New("invalid store type")
	ErrMissingStoreTypeAttr     = errors.New("missing store type")
	ErrMissingDirectoryPathAttr = errors.New("missing 'directory-path' attribute")
	ErrInvalidDirectoryPathAttr = errors.New("invalid 'directory-path' attribute (strings only)")
)
