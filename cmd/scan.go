package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"voltig/internal/models"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan system for all installed packages and versions",
	RunE: func(_cmd *cobra.Command, _args []string) error {
		pkgs, err := scanSystemPackages()
		if err != nil {
			return err
		}
		for _, p := range pkgs {
			logger.Info("Package found", "name", p.Name, "version", p.Version)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

// scanSystemPackages detects the OS/manager and returns all system packages
func scanSystemPackages() ([]models.PackageStatus, error) {
	osType := runtime.GOOS
	var cmd *exec.Cmd
	var parser func(string) []models.PackageStatus

	switch osType {
	case "darwin":
		cmd = exec.Command("brew", "list", "--versions")
		parser = parseBrew
	case "linux":
		// Try apt, fallback to pacman
		if _, err := exec.LookPath("dpkg-query"); err == nil {
			cmd = exec.Command("dpkg-query", "-W", "-f=${binary:Package}\t${Version}\n")
			parser = parseDpkg
		} else if _, err := exec.LookPath("pacman"); err == nil {
			cmd = exec.Command("pacman", "-Q")
			parser = parsePacman
		} else {
			return nil, fmt.Errorf("no supported package manager found for this OS")
		}
	case "windows":
		if _, err := exec.LookPath("choco"); err == nil {
			cmd = exec.Command("choco", "list", "--local-only")
			parser = parseChoco
		} else if _, err := exec.LookPath("winget"); err == nil {
			cmd = exec.Command("winget", "list")
			parser = parseWinget
		} else {
			return nil, fmt.Errorf("no supported package manager found for this OS")
		}
	default:
		logger.Error("Unsupported OS", "os", osType)
		return nil, fmt.Errorf("unsupported OS: %s", osType)
	}
	out, err := cmd.Output()
	if err != nil {
		logger.Error("Error detecting packages", "error", err)
		return nil, err
	}
	return parser(string(out)), nil
}

func parseBrew(out string) []models.PackageStatus {
	lines := strings.Split(out, "\n")
	var pkgs []models.PackageStatus
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) >= 2 {
			pkgs = append(pkgs, models.PackageStatus{Name: fields[0], Version: fields[1], Status: "installed"})
		}
	}
	return pkgs
}

func parseDpkg(out string) []models.PackageStatus {
	lines := strings.Split(out, "\n")
	var pkgs []models.PackageStatus
	for _, l := range lines {
		fields := strings.Split(l, "\t")
		if len(fields) >= 2 {
			pkgs = append(pkgs, models.PackageStatus{Name: fields[0], Version: fields[1], Status: "installed"})
		}
	}
	return pkgs
}

func parsePacman(out string) []models.PackageStatus {
	lines := strings.Split(out, "\n")
	var pkgs []models.PackageStatus
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) >= 2 {
			pkgs = append(pkgs, models.PackageStatus{Name: fields[0], Version: fields[1], Status: "installed"})
		}
	}
	return pkgs
}

func parseChoco(out string) []models.PackageStatus {
	lines := strings.Split(out, "\n")
	var pkgs []models.PackageStatus
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) >= 2 && !strings.HasPrefix(fields[0], "Chocolatey") {
			pkgs = append(pkgs, models.PackageStatus{Name: fields[0], Version: fields[1], Status: "installed"})
		}
	}
	return pkgs
}

func parseWinget(out string) []models.PackageStatus {
	lines := strings.Split(out, "\n")
	var pkgs []models.PackageStatus
	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) >= 3 && fields[0] != "Name" {
			pkgs = append(pkgs, models.PackageStatus{Name: fields[0], Version: fields[2], Status: "installed"})
		}
	}
	return pkgs
}
