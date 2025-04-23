package cmd

import (
	"os"
	"voltig/config"
	"voltig/internal/manager"
	"voltig/internal/models"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [pkg_name]...",
	Aliases: []string{"rm"},
	Short:   "Remove one or more packages",
	Args:    cobra.MinimumNArgs(1),
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

		var (
			failedRemovals  []string
			successRemovals []string
			notFound        []string
			targetPkgs      []config.Package
		)

		// Find packages to remove
		logger.Info("Removing specified packages", "packages", args)
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

		// Remove packages
		var pkgModels []models.Package
		for _, pkg := range targetPkgs {
			pkgModels = append(pkgModels, models.ToModel(pkg))
		}
		successRemovals, failedRemovals = manager.PkgOperation("Removing", "removed", pkgModels, m.Remove)

		// Print summary
		if len(successRemovals) > 0 {
			logger.Info("Successfully removed packages", "packages", successRemovals)
		}
		if len(failedRemovals) > 0 {
			logger.Error("Failed to remove packages", "packages", failedRemovals)
		}
		if len(notFound) > 0 {
			logger.Error("Packages not found in config", "packages", notFound)
		}

		// Exit with error code if any failures
		if len(failedRemovals) > 0 || len(notFound) > 0 {
			os.Exit(1)
		} else if len(successRemovals) > 0 {
			logger.Info("All requested packages removed successfully")
		} else {
			logger.Info("No packages were removed")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
