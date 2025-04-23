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

func main() {
	// Configure logger based on environment
	verbose := os.Getenv("VOLTIG_VERBOSE") == "true"
	if verbose {
		logger.SetLevel(log.DebugLevel)
		logger.Debug("Debug logging enabled")
	}

	// Log startup information
	logger.Info("Starting Voltig CLI",
		"version", "0.1.0",
		"verbose", verbose,
	)

	// Execute the root command
	cmd.Execute()
}
