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
	storeCmd.PersistentFlags().Int8VarP(&shareCount, "share-count", "a", 2,
		"share count (default: 2)")
	storeCmd.PersistentFlags().Int8VarP(&minSharesThreshold, "shares-threshold", "t", 2,
		"threshold of minimum shares which can be used retrieving file (default: 2)")
	storeCmd.PersistentFlags().BoolVarP(&encryptionTurnedOn, "encrypt", "e", false,
		"whether to encrypt file (AES) before storing it (default: false)")
	storeCmd.PersistentFlags().StringVarP(&encryptionSalt, "salt", "s", "",
		"32-byte long encryption salt (required if encryption mode is turned on)")

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

		gasper, err := pkg.NewGasper(extractStores(), &encryption.Settings{
			TurnedOn: decryptionTurnedOn,
			Salt:     decryptionSalt,
		})
		if err != nil {
			zap.L().Fatal("Failed to initialize Gasper", zap.Error(err))
		}

		zap.L().Info("Getting file shares")
		sharedFile, err := gasper.SharesFromFile(filePath, byte(shareCount), byte(minSharesThreshold))
		if err != nil {
			zap.L().Fatal("Failed to get file shares", zap.Error(err))
		}

		zap.L().Info("Check general stores availability")
		stores := gasper.Stores()
		availableStores := make([]storesPkg.Store, 0, len(stores))
		for _, store := range stores {
			if skip := checkStoreAvailability(store); skip {
				continue
			}

			availableStores = append(availableStores, store)
		}

		availableStoresCount := len(availableStores)
		if int(minSharesThreshold) > availableStoresCount {
			zap.L().Error("Not enough available stores", zap.Int8("Need", minSharesThreshold),
				zap.Int("Got", availableStoresCount), zap.Int8("Recommended", shareCount))
			return
		}

		zap.L().Info("Put shares in stores")
		i := 0
		for _, share := range sharedFile.Shares {
			if i > len(stores)-1 {
				zap.L().Error("All stores exhausted")
				return
			}

			store := stores[i]
			storeType := store.Type()
			i++

			zap.L().Debug("Available! Put share in store", zap.String("StoreType", storeType),
				zap.String("ShareID", share.ID))
			if err := store.Put(share); err != nil {
				zap.L().Error("Failed to put share in store", zap.String("StoreType", storeType),
					zap.Error(err))
				continue
			}
		}

		zap.L().Info("Success! Keep the following info for later use", zap.String("FileID", sharedFile.ID),
			zap.String("Checksum", sharedFile.Checksum))
	},
}
