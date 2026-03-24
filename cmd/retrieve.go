package cmd

import (
	"github.com/gasper/internal/encryption"
	"github.com/gasper/pkg"
	sharesPkg "github.com/gasper/pkg/shares"
	storesPkg "github.com/gasper/pkg/storage/stores"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	fileID             string
	destination        string
	checksum           string
	decryptionTurnedOn bool
	decryptionSalt     string
)

func init() {
	retrieveCmd.PersistentFlags().StringVarP(&fileID, "file-id", "i", "",
		"file id to retrieve (required)")
	retrieveCmd.PersistentFlags().StringVarP(&destination, "destination", "d", "",
		"where to save the retrieved file (required)")
	retrieveCmd.PersistentFlags().StringVarP(&checksum, "checksum", "m", "",
		"checksum of the shared file (required)")
	retrieveCmd.PersistentFlags().BoolVarP(&decryptionTurnedOn, "decrypt", "e", false,
		"whether file was encrypted before storing it (default: false)")
	retrieveCmd.PersistentFlags().StringVarP(&decryptionSalt, "salt", "s", "",
		"decryption salt (required if decryption mode is turned on)")

	if err := retrieveCmd.MarkPersistentFlagRequired("file-id"); err != nil {
		panic("Failed to mark 'file-id' flag as required")
	} else if err := retrieveCmd.MarkPersistentFlagRequired("destination"); err != nil {
		panic("Failed to mark 'destination' flag as required")
	} else if err := retrieveCmd.MarkPersistentFlagRequired("checksum"); err != nil {
		panic("Failed to mark 'checksum' flag as required")
	}

	rootCmd.AddCommand(retrieveCmd)
}

var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve a file",
	Long:  "Retrieve a file from the provided stores",
	Run: func(cmd *cobra.Command, args []string) {
		if decryptionTurnedOn && decryptionSalt == "" {
			zap.L().Fatal("Decryption salt is required when decryption mode is turned on")
		}

		gasper, err := pkg.NewGasper(extractStores(), &encryption.Settings{
			TurnedOn: decryptionTurnedOn,
			Salt:     decryptionSalt,
		})
		if err != nil {
			zap.L().Fatal("Failed to initialize Gasper", zap.Error(err))
		}

		sharedFile := &sharesPkg.SharedFile{
			ID:       fileID,
			Checksum: checksum,
			Shares:   make([]*sharesPkg.Share, 0),
		}

		zap.L().Info("Collect shares from stores")
		for _, store := range gasper.Stores() {
			store := store
			storeType := store.Type()

			if skip := checkStoreAvailability(store); skip {
				continue
			}

			zap.L().Debug("Available! Search share in store", zap.String("StoreType", storeType))
			share, err := store.Get(fileID)
			if err != nil {
				if err == storesPkg.ErrShareNotExists {
					zap.L().Debug("No match found in store, trying the next one", zap.String("StoreType",
						storeType))
					continue
				}

				zap.L().Error("Failed to search share in store", zap.String("StoreType", storeType),
					zap.Error(err))
				continue
			}

			sharedFile.Shares = append(sharedFile.Shares, share)
		}

		if len(sharedFile.Shares) == 0 {
			zap.L().Warn("No shares found for requested file ID", zap.String("FileID", fileID))
			return
		} else if len(sharedFile.Shares) < int(minSharesThreshold) {
			zap.L().Warn("Didn't find enough shares", zap.Int8("Need", minSharesThreshold),
				zap.Int("Got", len(sharedFile.Shares)))
			return
		}

		zap.L().Debug("Dump shared file")
		if err := gasper.DumpSharedFile(sharedFile, destination); err != nil {
			zap.L().Error("Failed dump shared file", zap.String("FileID", fileID),
				zap.String("Destination", destination), zap.Error(err))
			return
		}

		zap.L().Info("File retrieved successfully.", zap.String("FileID", fileID))
	},
}
