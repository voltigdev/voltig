// Package manager provides package management implementations for Voltig CLI.
package manager

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"voltig/internal/models"
	"voltig/pkg/logger"
)

// Allow stubbing exec.Command and exec.LookPath in tests
var execCommand = exec.Command
var execLookPath = exec.LookPath
var execOutput = func(cmd *exec.Cmd) ([]byte, error) { return cmd.Output() }

// BrewManager provides Homebrew package management on macOS.
type BrewManager struct{}
// Install package
func (b *BrewManager) Install(pkg models.Package, outputFn func(string)) error {
	// Handle multiple package names
	for _, name := range pkg.Name {
		logger.Info("Installing package", "name", name)
		
		args := []string{"install", name}
		if pkg.Version != "" && pkg.Version != "latest" {
			args = append(args, "--cask", fmt.Sprintf("%s@%s", name, pkg.Version))
		}
		cmd := execCommand("brew", args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		done := make(chan struct{}, 2)
		// Stream stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				lines := strings.Split(line, "\n")
				for _, l := range lines {
					trimmed := strings.TrimSpace(l)
					if trimmed != "" {
						if outputFn != nil {
							outputFn(trimmed)
						} else {
							logger.Info(trimmed)
						}
					}
				}
			}
			done <- struct{}{}
		}()
		// Stream stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				lines := strings.Split(line, "\n")
				for _, l := range lines {
					trimmed := strings.TrimSpace(l)
					if trimmed != "" {
						if outputFn != nil {
							outputFn(trimmed)
						} else {
							logger.Info(trimmed)
						}
					}
				}
			}
			done <- struct{}{}
		}()

		// Wait for both pipes to finish
		<-done
		<-done

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("failed to install %s: %w", name, err)
		}
	}
	
	return nil
}

// Update package
func (b *BrewManager) Update(pkg models.Package) error {
	// Handle multiple package names
	for _, name := range pkg.Name {
		logger.Info("Updating package", "name", name)
		
		cmd := execCommand("brew", "upgrade", name)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to update %s: %w", name, err)
		}
	}
	return nil
}

// Remove package
func (b *BrewManager) Remove(pkg models.Package, outputFn func(string)) error {
	// Handle multiple package names
	for _, name := range pkg.Name {
		logger.Info("Removing package", "name", name)
		
		cmd := execCommand("brew", "uninstall", name)
		cmd.Env = append(os.Environ(), "HOMEBREW_NO_AUTO_UPDATE=1")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}
		done := make(chan struct{}, 2)
		// Stream stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				lines := strings.Split(line, "\n")
				for _, l := range lines {
					trimmed := strings.TrimSpace(l)
					if trimmed != "" {
						// Only add prefix if line does not already start with [something]
						if !regexp.MustCompile(`^\[[^\]]+\] `).MatchString(trimmed) {
							trimmed = "[" + name + "] " + trimmed
						}
						if outputFn != nil {
							outputFn(trimmed)
						} else {
							logger.Info(trimmed)
						}
					}
				}
			}
			done <- struct{}{}
		}()
		// Stream stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				lines := strings.Split(line, "\n")
				for _, l := range lines {
					trimmed := strings.TrimSpace(l)
					if trimmed != "" {
						// Only add prefix if line does not already start with [something]
						if !regexp.MustCompile(`^\[[^\]]+\] `).MatchString(trimmed) {
							trimmed = "[" + name + "] " + trimmed
						}
						if outputFn != nil {
							outputFn(trimmed)
						} else {
							logger.Info(trimmed)
						}
					}
				}
			}
			done <- struct{}{}
		}()
		<-done
		<-done
		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("failed to remove %s: %w", name, err)
		}
	}
	return nil
}

// GetStatus checks the status of a package
func (b *BrewManager) GetStatus(pkg models.Package) (models.PackageStatus, error) {
	// For packages with multiple names, we'll check the first one
	if len(pkg.Name) == 0 {
		return models.PackageStatus{}, fmt.Errorf("package has no name")
	}
	
	name := pkg.Name[0]
	cmd := execCommand("brew", "list", "--versions", name)
	out, err := execOutput(cmd)
	if err == nil && len(out) > 0 {
		return models.PackageStatus{Name: name, Status: "installed", Version: string(out)}, nil
	}
	// Not managed by brew, check if binary exists in PATH
	if path, lookErr := execLookPath(name); lookErr == nil && path != "" {
		// Try to detect version manager
		var manager, version string
		if name == "node" || name == "nodejs" {
			if nvmPath, _ := execLookPath("nvm"); nvmPath != "" || isNvmNode(path) {
				manager = "nvm"
			}
		}
		if name == "python" || name == "python3" {
			if pyenvPath, _ := execLookPath("pyenv"); pyenvPath != "" || isPyenvPython(path) {
				manager = "pyenv"
			}
		}
		// Try to get version
		version = getBinaryVersion(name)
		status := "installed (external)"
		if manager != "" {
			status = "installed (" + manager + ")"
		}
		return models.PackageStatus{Name: name, Status: status, Version: version}, nil
	}
	return models.PackageStatus{Name: name, Status: "missing"}, nil
}

// isNvmNode checks if the node binary is managed by nvm
func isNvmNode(path string) bool {
	return strings.Contains(path, ".nvm/versions/")
}

// isPyenvPython checks if the python binary is managed by pyenv
func isPyenvPython(path string) bool {
	return strings.Contains(path, ".pyenv/versions/")
}

// getBinaryVersion runs '<binary> --version' and extracts the first version-like pattern (e.g., 1.2.3). Shows full output if ambiguous.
func getBinaryVersion(binary string) string {
	cmd := execCommand(binary, "--version")
	out, err := execOutput(cmd)
	if err != nil {
		return "unknown"
	}
	output := string(out)
	// Regex for version-like patterns (e.g., 1.2.3, v18.16.0, 25.0.0, etc.)
	re := regexp.MustCompile(`v?\d+\.\d+(\.\d+)?`)
	matches := re.FindAllString(output, -1)
	if len(matches) > 0 {
		return matches[0]
	}
	// fallback: return the entire output (trimmed)
	return strings.TrimSpace(output)
}

// IsAvailable checks if brew is available in the system PATH.
func (b *BrewManager) IsAvailable() bool {
	cmd := execCommand("brew", "--version")
	return cmd.Run() == nil
}