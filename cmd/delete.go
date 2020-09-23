package cmd

import (
	"github.com/gasper/internal/encryption"
	"github.com/gasper/pkg"
	storesPkg "github.com/gasper/pkg/storage/stores"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	deleteCmd.PersistentFlags().StringVarP(&fileID, "file-id", "i", "",
		"file id to retrieve (required)")

	if err := deleteCmd.MarkPersistentFlagRequired("file-id"); err != nil {
		panic("Failed to mark 'file-id' flag as required")
	}

	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a file",
	Long:  "Deletes a file from the provided stores",
	Run: func(cmd *cobra.Command, args []string) {
		gasper := pkg.NewGasper([]storesPkg.Store{}, &encryption.Settings{
			TurnedOn: decryptionTurnedOn,
			Salt:     decryptionSalt,
		})

		deletedShares := 0

		zap.L().Info("Delete shares from stores")
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

			zap.L().Debug("Available! Delete file from store", zap.String("StoreName", storeName))
			if err := store.Delete(fileID); err != nil {
				if err == pkg.ErrShareNotExists {
					zap.L().Debug("No match found in store, trying the next one", zap.String("StoreName",
						storeName))
					continue
				}

				zap.L().Error("Failed to delete share from store", zap.String("StoreName", storeName),
					zap.Error(err))
				continue // Best effort - keep trying other stores...
			}

			deletedShares++
		}

		if deletedShares == 0 {
			zap.L().Warn("No shares were found/deleted")
			return
		}

		zap.L().Info("File shares deleted successfully.", zap.String("FileID", fileID))
	},
}
