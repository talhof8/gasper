package stores

import "github.com/gasper/pkg/shares"

const TypeS3Store = "AmazonS3"

type AmazonS3Store struct {
	accessKey string
	secretKey string
}

func NewS3Store(accessKey string, secretKey string) (*AmazonS3Store, error){
	return &AmazonS3Store{
		accessKey: accessKey,
		secretKey: secretKey,
	}, nil
}

func (s AmazonS3Store) Type() string {
	return TypeS3Store
}

func (s AmazonS3Store) Available() (bool, error) {
	panic("implement me")
}

func (s AmazonS3Store) Put(share *shares.Share) error {
	panic("implement me")
}

func (s AmazonS3Store) Get(fileID string) (*shares.Share, error) {
	panic("implement me")
}

func (s AmazonS3Store) Delete(fileID string) error {
	panic("implement me")
}


