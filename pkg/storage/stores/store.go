package stores

import "github.com/gasper/pkg/shares"

// Store lets you store shares.
type Store interface {
	// Store type.
	Type() string

	// Is store available?
	// Useful especially for remote stores, such as ftp servers or s3 buckets.
	Available() (bool, error)

	// Puts a share in store.
	Put(share *shares.Share) error

	// Retrieves a share from store.
	// If no share with the given File ID exists, returns ErrShareNotExists.
	Get(fileID string) (*shares.Share, error)

	// Deletes a share from store.
	// If no share with the given File ID exists, returns ErrShareNotExists.
	Delete(fileID string) error
}
