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
		return AmazonS3Store{}
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
