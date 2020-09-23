package cmd

import (
	"github.com/gasper/internal/encryption"
	"github.com/gasper/pkg"
	storesPkg "github.com/gasper/pkg/storage/stores"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	filePath           string
	shareCount         int8
	minSharesThreshold int8
	encryptionTurnedOn bool
	encryptionSalt     string
)

func init() {
	storeCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "",
		"file to store (required)")
	storeCmd.PersistentFlags().Int8VarP(&shareCount, "share-count", "c", 3,
		"share count (default: 3)")
	storeCmd.PersistentFlags().Int8VarP(&minSharesThreshold, "shares-threshold", "t", 2,
		"threshold of minimum shares to use when restoring file (default: 2)")
	storeCmd.PersistentFlags().BoolVarP(&encryptionTurnedOn, "encrypt", "e", false,
		"whether to encrypt file (AES) before storing it (default: false)")
	storeCmd.PersistentFlags().StringVarP(&encryptionSalt, "salt", "s", "",
		"encryption salt (required if encryption mode is turned on)")

	if err := storeCmd.MarkPersistentFlagRequired("file"); err != nil {
		panic("Failed to mark 'file' flag as required")
	}

	rootCmd.AddCommand(storeCmd)
}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Store a file",
	Long:  "Store a file on the provided stores",
	Run: func(cmd *cobra.Command, args []string) {
		if minSharesThreshold > shareCount {
			zap.L().Fatal("Minimum shares threshold cannot be larger than share count")
		} else if encryptionTurnedOn && encryptionSalt == "" {
			zap.L().Fatal("Encryption salt is required when encryption mode is turned on")
		}

		gasper := pkg.NewGasper([]storesPkg.Store{}, &encryption.Settings{
			TurnedOn: encryptionTurnedOn,
			Salt:     encryptionSalt,
		})

		zap.L().Info("Getting file shares")
		sharedFile, err := gasper.SharesFromFile(filePath, byte(shareCount), byte(minSharesThreshold))
		if err != nil {
			zap.L().Fatal("Failed to get file shares", zap.Error(err))
		}

		zap.L().Info("Put shares in stores")
		stores := gasper.Stores()
		i := 0
		for _, share := range sharedFile.Shares {
			store := stores[i]
			storeName := store.Name()
			i++

			zap.L().Debug("Check store availability", zap.String("StoreName", storeName))
			available, err := store.Available()
			if err != nil {
				zap.L().Warn("Store availability check failed", zap.String("StoreName", storeName))
				continue
			} else if !available {
				zap.L().Debug("Skipping unavailable store", zap.String("StoreName", storeName))
				continue
			}

			zap.L().Debug("Available! Put share in store", zap.String("StoreName", storeName),
				zap.String("ShareID", share.ID))
			if err := store.Put(share); err != nil {
				zap.L().Error("Failed to put share in store", zap.String("StoreName", storeName),
					zap.Error(err))
				continue
			}
		}

		zap.L().Info("Success! Keep the following info for later use", zap.String("FileID", sharedFile.ID),
			zap.String("Checksum", sharedFile.Checksum))
	},
}
