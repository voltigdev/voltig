package cmd

import (
	"os"
	"voltig/config"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show loaded configuration",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			logger.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
		logger.Info("Loaded configuration", "config", cfg)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
