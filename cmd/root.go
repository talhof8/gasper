package cmd

import (
	"fmt"
	"github.com/gasper/internal/logging"
	storesPkg "github.com/gasper/pkg/storage/stores"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var (
	storesFile string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "gasper",
	Short: "Gasper lets you store files in a distributed manner on all sorts of different stores",
	Long: "Store your most sacred files by splitting and deploying them to different stores.\n" +
		"Retrieve them at any point, even if only a portion of the original destinations are available.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger, err := logging.NewLogger("gasper", verbose)
		if err != nil {
			return err
		}

		zap.ReplaceGlobals(logger)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		_ = zap.L().Sync()
		_ = zap.S().Sync()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&storesFile, "stores-config", "c", "", "stores config file (required)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "extra verbosity")

	if err := rootCmd.MarkPersistentFlagRequired("stores-config"); err != nil {
		panic("Failed to mark 'stores-config' flag as required")
	}
}

func extractStores() []storesPkg.Store {
	config := viper.New()
	config.SetConfigFile(storesFile)
	if err := config.ReadInConfig(); err != nil {
		zap.L().Fatal("Failed to read stores config file", zap.String("Path", filePath), zap.Error(err))
	}

	storesConfigRaw := config.Get("stores")
	if storesConfigRaw == nil {
		zap.L().Fatal("Empty stores config")
	}

	storesConfig, ok := storesConfigRaw.([]interface{})
	if !ok {
		zap.L().Fatal("Invalid configuration scheme")
	}

	stores := make([]storesPkg.Store, 0)
	for _, storeConfig := range storesConfig {
		storeConfigMap, ok := storeConfig.(map[string]interface{})
		if !ok {
			zap.L().Fatal("Invalid configuration scheme")
		}

		store, err := storesPkg.FromConfig(storeConfigMap)
		if err != nil {
			zap.L().Fatal("Failed to create store from config", zap.Any("RawConfig", storeConfig),
				zap.Error(err))
		}

		stores = append(stores, store)
	}
	return stores
}

func checkStoreAvailability(store storesPkg.Store) bool {
	storeName := store.Name()

	zap.L().Debug("Check store availability", zap.String("StoreName", storeName))
	available, err := store.Available()
	if err != nil {
		zap.L().Warn("Store availability check failed", zap.String("StoreName", storeName), zap.Error(err))
		return true
	} else if !available {
		zap.L().Debug("Skipping unavailable store", zap.String("StoreName", storeName))
		return true
	}

	return false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
