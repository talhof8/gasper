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
	retrieveCmd.PersistentFlags().StringVarP(&checksum, "checksum", "c", "",
		"checksum of the shared file (required)")
	retrieveCmd.PersistentFlags().BoolVarP(&decryptionTurnedOn, "encrypt", "e", false,
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
	Short: "Retrieves a file",
	Long:  "Retrieves a file from the provided stores",
	Run: func(cmd *cobra.Command, args []string) {
		if decryptionTurnedOn && decryptionSalt == "" {
			zap.L().Fatal("Decryption salt is required when decryption mode is turned on")
		}

		gasper := pkg.NewGasper([]storesPkg.Store{}, &encryption.Settings{
			TurnedOn: decryptionTurnedOn,
			Salt:     decryptionSalt,
		})

		sharedFile := &sharesPkg.SharedFile{
			ID:       fileID,
			Checksum: checksum,
			Shares:   make([]*sharesPkg.Share, 0),
		}

		zap.L().Info("Collect shares from stores")
		for _, store := range gasper.Stores() {
			store := store
			storeName := store.Name()

			zap.L().Debug("Check store availability", zap.String("StoreName", storeName))
			available, err := store.Available()
			if err != nil {
				zap.L().Warn("Store availability check failed", zap.String("StoreName", storeName))
				continue
			} else if !available {
				zap.L().Debug("Skipping unavailable store", zap.String("StoreName", storeName))
				continue
			}

			zap.L().Debug("Available! Search share in store", zap.String("StoreName", storeName))
			share, err := store.Get(fileID)
			if err != nil {
				if err == pkg.ErrShareNotExists {
					zap.L().Debug("No match found in store, trying the next one", zap.String("StoreName",
						storeName))
					continue
				}

				zap.L().Error("Failed to search share in store", zap.String("StoreName", storeName),
					zap.Error(err))
				continue
			}

			sharedFile.Shares = append(sharedFile.Shares, share)
		}

		if len(sharedFile.Shares) == 0 {
			zap.L().Warn("No shares found for requested file ID", zap.String("FileID", fileID))
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
