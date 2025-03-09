package ui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// GetMenuTheme returns the styled theme for menu UI
func GetMenuTheme() *huh.Theme {
	t := huh.ThemeBase()

	// Selected option
	t.Focused.SelectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color("#FD0053"))

	// Unselected options
	t.Focused.Option = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FD0053"))

	// Menu title
	t.Focused.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FD0053")).
		BorderBottom(true).
		BorderStyle(lipgloss.RoundedBorder())

	t.Focused.SelectSelector = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FD0053")).
		SetString("> ")

	t.Focused.Base = lipgloss.NewStyle().
		Padding(1, 1, 0, 1).
		Foreground(lipgloss.Color("#00BDED"))

	return t
}
