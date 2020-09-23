package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/codahale/sss"
	"github.com/gasper/internal/encryption"
	sharesPkg "github.com/gasper/pkg/shares"
	storesPkg "github.com/gasper/pkg/storage/stores"
	"github.com/lithammer/shortuuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"strconv"
)

// todo: more elegant and efficient way to read & write big files.

// Gasper lets you store, load, and delete files in a multi-part, distributed manner, using on Shamir's Secret Sharing.
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

// Retrieves stores.
func (g *Gasper) Stores() []storesPkg.Store {
	return g.stores
}

// Splits file into its shares.
func (g *Gasper) SharesFromFile(filePath string, shareCount, minSharesThreshold byte) (*sharesPkg.SharedFile, error) {
	if minSharesThreshold > shareCount {
		return nil, ErrInvalidSharesThreshold
	}

	fileID := shortuuid.New()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "read file '%s'", filePath)
	}

	encryptedData, err := g.encryptor.Encrypt(data)
	if err != nil {
		return nil, errors.WithMessage(err, "encrypt data")
	}

	sharesBytes, err := sss.Split(shareCount, minSharesThreshold, encryptedData)
	if err != nil {
		return nil, errors.WithMessage(err, "split data to shares")
	}

	shares := make([]*sharesPkg.Share, 0, len(sharesBytes))
	for shareID, shareBytes := range sharesBytes {
		share := &sharesPkg.Share{
			ID:     strconv.Itoa(int(shareID)),
			FileID: fileID,
			Data:   shareBytes,
		}

		shares = append(shares, share)
	}

	checksum := md5.Sum(data)

	return &sharesPkg.SharedFile{
		ID:       fileID,
		Checksum: hex.EncodeToString(checksum[:]),
		Shares:   shares,
	}, nil
}

// Dumps shared file to a local filesystem destination.
// If md5 checksum is set, will use it to check file authenticity, otherwise will skip checksum check.
func (g *Gasper) DumpSharedFile(sharedFile *sharesPkg.SharedFile, destination string) error {
	if sharedFile == nil {
		return ErrNilSharedFile
	}

	rawShares := make(map[byte][]byte, len(sharedFile.Shares))
	for _, share := range sharedFile.Shares {
		shareIDInt, err := strconv.Atoi(share.ID)
		if err != nil {
			return errors.WithMessage(err, "convert share ID from string to int")
		}
		rawShares[byte(shareIDInt)] = share.Data
	}

	combinedBytes := sss.Combine(rawShares)

	decryptedData, err := g.encryptor.Decrypt(combinedBytes)
	if err != nil {
		return errors.WithMessage(err, "decrypt data")
	}

	if sharedFile.Checksum != "" {
		if err := g.validateChecksum(decryptedData, sharedFile.Checksum); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(destination, decryptedData, os.ModePerm); err != nil {
		return errors.WithMessagef(err, "write file '%s'", destination)
	}
	return nil
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
