package cmd

import (
	"fmt"
	"github.com/gasper/internal/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "gasper",
	Short: "Gasper lets you backup files in a distributed manner using Shamir's Secret Sharing",
	Long: "Backup your most sacred files by splitting and deploying them on different destinations.\n" +
		"Retrieve them at any point, given only a portion of the original destinations being available.",
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "extra verbosity")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
