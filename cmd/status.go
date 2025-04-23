package cmd

import (
	"os"
	"voltig/config"
	"voltig/internal/manager"
	"voltig/internal/models"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Show installation status of all packages",
	Args:    cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			logger.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
		m := manager.ForOS()
		if m == nil {
			logger.Error("No supported package manager found for this OS")
			os.Exit(1)
		}
		for _, pkg := range cfg.Packages {
			status, _ := m.GetStatus(models.ToModel(pkg))
			logger.Info("Package status", "name", status.Name, "status", status.Status, "version", status.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
