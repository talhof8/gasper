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
	Short: "Delete a file",
	Long:  "Delete a file from the provided stores",
	Run: func(cmd *cobra.Command, args []string) {
		gasper, err := pkg.NewGasper(extractStores(), &encryption.Settings{
			TurnedOn: false,
		})
		if err != nil {
			zap.L().Fatal("Failed to initialize Gasper", zap.Error(err))
		}

		deletedShares := 0

		zap.L().Info("Delete shares from stores")
		for _, store := range gasper.Stores() {
			store := store
			storeType := store.Type()

			if skip := checkStoreAvailability(store); skip {
				continue
			}

			zap.L().Debug("Available! Delete file from store", zap.String("StoreType", storeType))
			if err := store.Delete(fileID); err != nil {
				if err == storesPkg.ErrShareNotExists {
					zap.L().Debug("No match found in store, trying the next one", zap.String("StoreType",
						storeType))
					continue
				}

				zap.L().Error("Failed to delete share from store", zap.String("StoreType", storeType),
					zap.Error(err))
				continue // Best effort - keep trying other stores...
			}

			deletedShares++
		}

		if deletedShares == 0 {
			zap.L().Warn("No shares were found/deleted")
			return
		}

		zap.L().Info("File shares deleted successfully.", zap.String("FileID", fileID),
			zap.Int("DeletedShares", deletedShares))
	},
}
