package cmd

import (
	"os"
	"os/exec"
	"voltig/config"
	"voltig/pkg/logger"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "voltig",
	Short: "A cross-platform package manager CLI",
	Long:  `A cross-platform CLI for managing system packages with TUI support.`,
	// Enable command grouping
	GroupID: "main",
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "completion [shell]",
		Short: "Generate shell completion script (bash|zsh|fish)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				if err := rootCmd.GenBashCompletion(os.Stdout); err != nil {
					logger.Error("Failed to generate bash completion", "error", err)
				}
			case "zsh":
				if err := rootCmd.GenZshCompletion(os.Stdout); err != nil {
					logger.Error("Failed to generate zsh completion", "error", err)
				}
			case "fish":
				if err := rootCmd.GenFishCompletion(os.Stdout, true); err != nil {
					logger.Error("Failed to generate fish completion", "error", err)
				}
			default:
				logger.Error("Unknown shell", "shell", args[0])
				os.Exit(1)
			}
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "lint",
		Short: "Lint and validate your voltig.yml config file",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				logger.Error("Config error", "error", err)
				os.Exit(1)
			}
			// Check for duplicate package names
			nameSet := make(map[string]struct{})
			for _, pkg := range cfg.Packages {
				for _, name := range pkg.Name {
					if _, exists := nameSet[name]; exists {
						logger.Error("Duplicate package name", "name", name)
						os.Exit(1)
					}
					nameSet[name] = struct{}{}
				}
				if len(pkg.Name) == 0 || pkg.Manager == "" {
					logger.Error("Missing required fields in package", "package", pkg)
					os.Exit(1)
				}
			}
			// Check for protected command overrides
			protected := map[string]struct{}{"install": {}, "update": {}, "remove": {}, "status": {}, "tui": {}, "help": {}, "completion": {}, "lint": {}}
			for name := range cfg.Commands {
				if _, found := protected[name]; found {
					logger.Error("Protected command cannot be overridden", "command", name)
					os.Exit(1)
				}
			}
			logger.Info("Config validation successful")
		},
	})
}

func Execute() {
	// List of protected/core commands
	protected := map[string]struct{}{
		"install": {}, "update": {}, "remove": {}, "status": {}, "tui": {}, "help": {}, "completion": {}, "lint": {}, "config": {}, "scan": {},
	}

	// Assign core commands to their group
	for _, cmd := range rootCmd.Commands() {
		// Skip help command
		if cmd.Name() == "help" {
			continue
		}

		// Assign utility commands
		if cmd.Name() == "completion" || cmd.Name() == "lint" || cmd.Name() == "config" {
			cmd.GroupID = "utility"
			continue
		}

		// Assign all other built-in commands to core group
		if _, found := protected[cmd.Name()]; found {
			cmd.GroupID = "core"
		}
	}
	// Load config to register user commands
	cfg, err := config.LoadConfig(configFile)
	if err == nil && cfg.Commands != nil {
		for name, c := range cfg.Commands {
			if _, found := protected[name]; found {
				logger.Error("User command not allowed", "command", name, "reason", "name is protected by core CLI")
				os.Exit(1)
			}
			cmdName := name
			cmdDef := c
			cmd := &cobra.Command{
				Use:     cmdName,
				Short:   "✨ " + cmdDef.Summary,
				GroupID: "project",
				Run: func(cmd *cobra.Command, args []string) {
					logger.Info(cmdDef.Summary)
					shellCmd := exec.Command("sh", "-c", cmdDef.Command)
					shellCmd.Stdout = os.Stdout
					shellCmd.Stderr = os.Stderr
					shellCmd.Stdin = os.Stdin
					if err := shellCmd.Run(); err != nil {
						logger.Error("Command failed", "error", err)
						os.Exit(1)
					}
				},
			}
			rootCmd.AddCommand(cmd)
		}
	}
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "voltig.yml", "Path to YAML config file")

	// Set up command groups
	rootCmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core Package Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "project",
		Title: "Project Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "utility",
		Title: "Utility Commands:",
	})
}
