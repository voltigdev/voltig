package cmd

import "github.com/charmbracelet/lipgloss"

var (
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFAF")).
			Padding(0, 1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF005F")).
			Bold(true)
)
