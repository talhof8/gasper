package stores

import (
	"fmt"
	"github.com/gasper/pkg/shares"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const TypeLocalStore = "local"

// Stores files in a local directory,
// Note: needs to get an absolute path.
type LocalStore struct {
	directoryPath string
}

func NewLocalStore(directoryPath string) (*LocalStore, error) {
	return &LocalStore{
		directoryPath: directoryPath,
	}, nil
}

func (ls *LocalStore) Name() string {
	return "local-store"
}

func (ls *LocalStore) Available() (bool, error) {
	_, err := os.Stat(ls.directoryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (ls *LocalStore) Put(share *shares.Share) error {
	filePath := path.Join(ls.directoryPath, ls.filename(share))
	return ioutil.WriteFile(filePath, share.Data, os.ModePerm)
}

func (ls *LocalStore) Get(fileID string) (*shares.Share, error) {
	filePath, err := ls.findFileByID(fileID)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "read file '%s'", filePath)
	}

	splitFilename := strings.Split(path.Base(filePath), ".")
	if len(splitFilename) != 3 { // Note: should not happen, otherwise glob wouldn't match. But just in case...
		return nil, errors.New("invalid file format (should be: '<file-id>.<share-id>.gasper')")
	}

	shareID := splitFilename[1]

	return &shares.Share{
		ID:     shareID,
		FileID: fileID,
		Data:   data,
	}, nil
}

func (ls *LocalStore) Delete(fileID string) error {
	filePath, err := ls.findFileByID(fileID)
	if err != nil {
		return err
	}

	return os.RemoveAll(filePath)
}

func (ls *LocalStore) filename(share *shares.Share) string {
	return fmt.Sprintf("%s.%s.gasper", share.FileID, share.ID)
}

func (ls *LocalStore) findFileByID(fileID string) (string, error) {
	pattern := path.Join(ls.directoryPath, fmt.Sprintf("%s.*.gasper", fileID))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", errors.WithMessagef(err, "glob pattern '%s'", pattern)
	}

	if len(matches) == 0 {
		return "", ErrShareNotExists
	} else if len(matches) > 1 {
		return "", ErrMoreThanOneMatch
	}

	return matches[0], nil
}
