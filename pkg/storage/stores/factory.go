package stores

func FromConfig(config map[string]interface{}) (Store, error) {
	storeType, ok := config["type"]
	if !ok {
		return nil, ErrMissingStoreTypeAttr
	}

	switch storeType {
	case TypeLocalStore:
		return localStore(config)
	case TypeS3Store:
		return amazonS3Store(config)
	}

	return nil, ErrInvalidStoreType
}

func localStore(config map[string]interface{}) (Store, error) {
	directoryPathRaw, ok := config["directory-path"]
	if !ok {
		return nil, ErrMissingDirectoryPathAttr
	}

	directoryPath, ok := directoryPathRaw.(string)
	if !ok {
		return nil, ErrInvalidDirectoryPathAttr
	}

	return NewLocalStore(directoryPath)
}

func amazonS3Store(config map[string]interface{}) (Store, error) {
	s3AccessKeyRaw, ok := config["access-key"]
	if !ok {
		return nil, ErrMissingAmazonS3AccessKeyAttr
	}

	s3AccessKey, ok := s3AccessKeyRaw.(string)
	if !ok {
		return nil, ErrInvalidAmazonS3AccessKeyAttr
	}

	s3SecretKeyRaw, ok := config["secret-key"]
	if !ok {
		return nil, ErrMissingAmazonS3SecretKeyAttr
	}

	s3SecretKey, ok := s3SecretKeyRaw.(string)
	if !ok {
		return nil, ErrInvalidAmazonS3SecretAttr
	}

	return NewS3Store(s3AccessKey, s3SecretKey)
}
