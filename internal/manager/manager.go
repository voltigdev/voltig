package manager

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"voltig/internal/models"
	"voltig/pkg/logger"
)

// PackageManager defines the interface for a system package manager used by Voltig CLI.
type PackageManager interface {
	Install(pkg models.Package, outputFn func(string)) error
	Update(pkg models.Package) error
	Remove(pkg models.Package, outputFn func(string)) error
	GetStatus(pkg models.Package) (models.PackageStatus, error)
	IsAvailable() bool
}

// ForOS returns the appropriate PackageManager for the current OS.
func ForOS() PackageManager {
	if runtime.GOOS == "darwin" {
		bm := &BrewManager{}
		if bm.IsAvailable() {
			return bm
		}
	}
	return nil
}

// PkgOperation performs an operation (install, remove, etc.) on a list of packages using the provided opFunc.
func PkgOperation(opName, opPast string, pkgs []models.Package, opFunc func(models.Package, func(string)) error) (successes, failures []string) {
	total := len(pkgs)
	for i, pkg := range pkgs {
		// For packages with multiple names, we'll log all names
		pkgNameStr := strings.Join(pkg.Name, ", ")
		logger.Info(opName+" package", "package", pkgNameStr, "progress", fmt.Sprintf("%d/%d", i+1, total))

		showOutput := func(line string) {
			lines := strings.Split(line, "\n")
			last := ""
			for _, l := range lines {
				trimmed := strings.TrimSpace(l)
				if trimmed != "" {
					last = trimmed
				}
			}
			if last != "" {
				if _, err := fmt.Fprintf(os.Stdout, "\r[%s] %s", pkgNameStr, last); err != nil {
					logger.Error("Failed to write to stdout", "error", err)
				}
			}
		}
		err := opFunc(pkg, showOutput)

		if _, err := fmt.Fprint(os.Stdout, "\n"); err != nil {
			logger.Error("Failed to write newline to stdout", "error", err)
		}

		if err != nil {
			logger.Error("Failed to "+opName+" package", "package", pkgNameStr, "error", err)
			// Add each individual package name to the failures list
			failures = append(failures, pkg.Name...)
		} else {
			logger.Info("Successfully "+opPast+" package", "package", pkgNameStr)
			// Add each individual package name to the successes list
			successes = append(successes, pkg.Name...)
		}
	}
	return successes, failures
}