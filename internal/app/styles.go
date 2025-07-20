package app

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	mainMenuCurosrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#25A065")) // green

	configureGroupStatusStyleEditing    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Background(lipgloss.Color("235")).Padding(0, 1)
	configureGroupStatusStyleNavigation = lipgloss.NewStyle().Foreground(lipgloss.Color("#436f94")).Background(lipgloss.Color("235")).Padding(0, 1).Bold(true)

	errorMessageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true) // red
	infoMessageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true) // green
	warningMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true) // orange
)
