package cmd

import (
	"fmt"
	"github.com/gasper/internal/logging"
	"github.com/spf13/cobra"
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
	rootCmd.PersistentFlags().StringVarP(&storesFile, "stores-config", "s", "", "stores configuraion file (required)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "extra verbosity")

	if err := rootCmd.MarkPersistentFlagRequired("stores-config"); err != nil {
		panic("Failed to mark 'stores-config' flag as required")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
