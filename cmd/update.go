package cmd

import (
	"os"
	"voltig/config"
	"voltig/internal/manager"
	"voltig/internal/models"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Aliases: []string{"u"},
	Use:   "update [pkg_name...]",
	Short: "Update specific or all packages",
	Args:  cobra.ArbitraryArgs,
	Run: func(_ *cobra.Command, args []string) {
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
		var targetPkgs []config.Package
		var notFound []string
		// Determine which packages to update
		if len(args) == 0 {
			// Update all packages from config
			targetPkgs = cfg.Packages
			logger.Info("Updating all packages from config", "count", len(targetPkgs))
		} else {
			// Update only specified packages
			logger.Info("Updating specified packages", "packages", args)
			for _, arg := range args {
				found := false
				for _, pkg := range cfg.Packages {
					// Check if the package name is in the list of names
					for _, name := range pkg.Name {
						if name == arg {
							targetPkgs = append(targetPkgs, pkg)
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					logger.Error("Package not found in config", "package", arg)
					notFound = append(notFound, arg)
				}
			}
		}
		// Update the target packages
		for _, pkg := range targetPkgs {
			_ = m.Update(models.ToModel(pkg))
		}
		if len(notFound) > 0 {
			logger.Error("Failed to update packages", "packages", notFound)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
