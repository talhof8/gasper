package stores

import "github.com/pkg/errors"

var (
	ErrShareNotExists   = errors.New("share doesn't exist in store")
	ErrMoreThanOneMatch = errors.New("found more than one match for share")

	// Missing/invalid config attributes errors.
	ErrInvalidStoreType             = errors.New("invalid store type")
	ErrMissingStoreTypeAttr         = errors.New("missing store type")
	ErrMissingDirectoryPathAttr     = errors.New("missing 'directory-path' attribute")
	ErrInvalidDirectoryPathAttr     = errors.New("invalid 'directory-path' attribute (strings only)")
	ErrMissingAmazonS3AccessKeyAttr = errors.New("missing 'access-key' attribute")
	ErrInvalidAmazonS3AccessKeyAttr = errors.New("invalid 'access-key' attribute (strings only)")
	ErrMissingAmazonS3SecretKeyAttr = errors.New("missing 'secret-key' attribute")
	ErrInvalidAmazonS3SecretKeyAttr = errors.New("invalid 'secret-key' attribute (strings only)")
	ErrMissingAmazonS3RegionAttr = errors.New("missing 'region' attribute")
	ErrInvalidAmazonS3RegionAttr = errors.New("invalid 'region' attribute (strings only)")
)
