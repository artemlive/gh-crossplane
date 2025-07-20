package app

import "github.com/charmbracelet/lipgloss"

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	infoMessageStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true) // orange
	mainMenuCurosrStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#25A065"))            // green
	mainMenuSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5"))

	configureGroupStatusStyleEditing    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Background(lipgloss.Color("235")).Padding(0, 1)
	configureGroupStatusStyleNavigation = lipgloss.NewStyle().Foreground(lipgloss.Color("#436f94")).Background(lipgloss.Color("235")).Padding(0, 1).Bold(true)
)
