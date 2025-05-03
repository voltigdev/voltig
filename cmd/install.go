package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"voltig/config"
	"voltig/internal/manager"
	"voltig/internal/models"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:     "install [pkg_name...]",
	Aliases: []string{"i"},
	Short:   "Install specific or all packages",
	Args:    cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, args []string) {
		fmt.Println(HeaderStyle.Render("🔧 Voltig: Installing Packages"))

		if err := ensureHomebrew(); err != nil {
			fmt.Println(ErrorStyle.Render("Failed to install Homebrew:", err.Error()))
			os.Exit(1)
		}

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
			failedInstalls  []string
			successInstalls []string
			notFound        []string
			targetPkgs      []config.Package
		)

		// Determine which packages to install
		if len(args) == 0 {
			// Install all packages from config
			targetPkgs = cfg.Packages
			logger.Info("Installing all packages from config", "count", len(targetPkgs))
		} else {
			// Install only specified packages
			logger.Info("Installing specified packages", "packages", args)
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

		// Install packages
		var pkgModels []models.Package
		for _, pkg := range targetPkgs {
			pkgModels = append(pkgModels, models.ToModel(pkg))
		}
		successInstalls, failedInstalls = manager.PkgOperation("Installing", "installed", pkgModels, m.Install)

		// Print summary
		if len(successInstalls) > 0 {
			logger.Info("Successfully installed packages", "packages", successInstalls)
		}
		if len(failedInstalls) > 0 {
			logger.Error("Failed to install packages", "packages", failedInstalls)
		}
		if len(notFound) > 0 {
			logger.Error("Packages not found in config", "packages", notFound)
		}

		// Exit with error code if any failures
		if len(failedInstalls) > 0 || len(notFound) > 0 {
			os.Exit(1)
		} else if len(successInstalls) > 0 {
			logger.Info("All requested packages installed successfully")
		} else {
			logger.Info("No packages were installed")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func ensureHomebrew() error {
	_, err := exec.LookPath("brew")
	if err == nil {
		return nil
	}
	fmt.Println(HeaderStyle.Render("Homebrew not found. Installing Homebrew..."))
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		return fmt.Errorf("automatic Homebrew installation is not supported on this OS")
	}
}
