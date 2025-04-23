// Package util provides utility functions for formatting and pretty-printing output.
package util

import (
	"github.com/charmbracelet/lipgloss"
)

var brewStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250"))

var pacmanStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250"))

// StyleLine returns a prettily styled output line for CLI or TUI
func StyleLine(line string) string {
	switch line {
	case "brew":
		return brewStyle.Render(	line)
	case "pacman":
		return pacmanStyle.Render(line)
	default:
		return line
	}
}

