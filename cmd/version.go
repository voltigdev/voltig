package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Set these at build time using -ldflags
	version   = "dev"
	commit    = "none"
	date      = "unknown"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show the version of voltig CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("voltig version %s\ncommit: %s\nbuilt at: %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
