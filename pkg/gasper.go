package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/codahale/sss"
	"github.com/gasper/internal/encryption"
	"github.com/gasper/internal/shares"
	"github.com/gasper/pkg/storage"
	storesPkg "github.com/gasper/pkg/storage/stores"
	"github.com/lithammer/shortuuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strconv"
)

// todo: more elegant and efficient way to read & write big files.

// Gasper lets you store, load, and delete files in a multi-part, distributed manner, based on Shamir's Secret Sharing.
// It holds a list of stores being used for distribution, and encryption settings.
type Gasper struct {
	stores    []storesPkg.Store
	encryptor *encryption.Encryptor
}

func NewGasper(stores []storesPkg.Store, encryptionSettings *encryption.Settings) *Gasper {
	return &Gasper{
		stores:    stores,
		encryptor: encryption.NewEncryptor(encryptionSettings),
	}
}

// Distributes a file across stores.
// Retrieves a unique file id which can should be used for restore, a md5 checksum of the original file, and
// an error if one occurred.
func (g *Gasper) Distribute(filePath string, shareCount, minSharesThreshold byte) (string, string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", errors.WithMessagef(err, "read file '%s'", filePath)
	}

	encryptedData, err := g.encryptor.Encrypt(data)
	if err != nil {
		return "", "", errors.WithMessage(err, "encrypt data")
	}

	sharesBytes, err := sss.Split(shareCount, minSharesThreshold, encryptedData)
	if err != nil {
		return "", "", errors.WithMessage(err, "split data to shares")
	}

	fileID, err := g.distributeShares(sharesBytes)
	if err != nil {
		return "", "", err
	}

	checksum := md5.Sum(data)
	return fileID, hex.EncodeToString(checksum[:]), nil
}

// Restores a file given its ID to a destination file.
// If md5 checksum is provided, will use it to check file authenticity, otherwise will skip checksum check.
func (g *Gasper) Restore(fileID string, destination string, originalChecksum string) error {
	rawShares, err := g.collectShares(fileID)
	if err != nil {
		return err
	}

	combinedBytes := sss.Combine(rawShares)

	decryptedData, err := g.encryptor.Decrypt(combinedBytes)
	if err != nil {
		return errors.WithMessage(err, "decrypt data")
	}

	if originalChecksum != "" {
		if err := g.validateChecksum(decryptedData, originalChecksum); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(destination, decryptedData, os.ModePerm); err != nil {
		return errors.WithMessagef(err, "write file '%s'", destination)
	}
	return nil
}

func (g *Gasper) distributeShares(sharesBytes map[byte][]byte) (string, error) {
	fileID := shortuuid.New()

	i := 0

	for shareID, shareBytes := range sharesBytes {
		share := &shares.Share{
			ID:     strconv.Itoa(int(shareID)),
			FileID: fileID,
			Data:   shareBytes,
		}

		store := g.stores[i]

		available, err := store.Available()
		if err != nil {
			zap.L().Warn("Store availability check failed", zap.String("StoreName", store.Name()))
			continue
		} else if !available {
			zap.L().Debug("Skipping unavailable store", zap.String("StoreName", store.Name()))
			continue
		}

		if err := store.Put(share); err != nil {
			return "", errors.WithMessagef(err, "put share in store '%s'", store.Name())
		}

		i++
	}

	return fileID, nil
}

func (g *Gasper) collectShares(fileID string) (map[byte][]byte, error) {
	rawShares := make(map[byte][]byte, 0)
	for _, store := range g.stores {
		available, err := store.Available()
		if err != nil {
			zap.L().Warn("Store availability check failed", zap.String("StoreName", store.Name()))
			continue
		} else if !available {
			zap.L().Debug("Skipping unavailable store", zap.String("StoreName", store.Name()))
			continue
		}

		share, err := store.Get(fileID)
		if err != nil {
			if err == storage.ErrShareNotExists {
				zap.L().Debug("Share for file doesn't exist in this store, trying the next one",
					zap.String("StoreName", store.Name()))
				continue
			}

			return nil, errors.WithMessagef(err, "get share from store '%s'", store.Name())
		}

		shareIDInt, err := strconv.Atoi(share.ID)
		if err != nil {
			return nil, errors.WithMessage(err, "convert share ID from string to int")
		}
		rawShares[byte(shareIDInt)] = share.Data
	}

	return rawShares, nil
}

func (g *Gasper) validateChecksum(decryptedData []byte, originalChecksum string) error {
	currentChecksumBytes := md5.Sum(decryptedData)
	currentChecksum := hex.EncodeToString(currentChecksumBytes[:])

	if currentChecksum != originalChecksum {
		return errors.Errorf("got corrupt data, checksums didn't match (original: '%s', got: '%s')",
			originalChecksum, currentChecksum)
	}
	return nil
}
