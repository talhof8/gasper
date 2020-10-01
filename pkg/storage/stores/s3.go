package stores

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gasper/pkg/shares"
	"github.com/pkg/errors"
)

const TypeS3Store = "AmazonS3"

type AmazonS3Store struct {
	accessKey string
	secretKey string
	region    string
}

func NewS3Store(accessKey, secretKey, region string) (*AmazonS3Store, error){
	return &AmazonS3Store{
		accessKey: accessKey,
		secretKey: secretKey,
		region: region,
	}, nil
}

func (s AmazonS3Store) Type() string {
	return TypeS3Store
}

func (s AmazonS3Store) Available() (bool, error) {
	if len(s.accessKey) == 0 || len(s.secretKey) == 0 {
		return false, nil
	}
	return true, nil
}

func (s AmazonS3Store) Put(share *shares.Share) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.region),
		Credentials: credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     s.accessKey,
			SecretAccessKey: s.secretKey,
		}),
	})
	if err != nil {
		// TODO: Add reason
		return errors.WithMessagef(err, "Unable to upload file to S3")
	}

}

func (s AmazonS3Store) Get(fileID string) (*shares.Share, error) {
	panic("implement me")
}

func (s AmazonS3Store) Delete(fileID string) error {
	panic("implement me")
}


