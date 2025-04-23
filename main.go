/*
Package main provides the entry point for Voltig CLI.
*/
package main

import (
	"os"
	"voltig/cmd"
	"voltig/pkg/logger"

	"github.com/charmbracelet/log"
)

// Note: version, commit, and date variables are injected via -ldflags at build time for the version command.
func main() {
	// Configure logger based on environment
	verbose := os.Getenv("VOLTIG_VERBOSE") == "false"
	if verbose {
		logger.SetLevel(log.DebugLevel)
		logger.Debug("Debug logging enabled")
	}

	// Execute the root command
	cmd.Execute()
}
