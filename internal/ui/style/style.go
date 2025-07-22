package style

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	MainMenuCurosrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#25A065")) // green

	ConfigureGroupStatusStyleEditing    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Background(lipgloss.Color("235")).Padding(0, 1)
	ConfigureGroupStatusStyleNavigation = lipgloss.NewStyle().Foreground(lipgloss.Color("#436f94")).Background(lipgloss.Color("235")).Padding(0, 1).Bold(true)

	ErrorMessageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true) // red
	InfoMessageStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true) // green
	WarningMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true) // orange

	InactiveTabStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(lipgloss.Color("240"))
	ActiveTabStyle   = InactiveTabStyle.BorderBottom(false).Bold(true)

	RepoPreviewStyle = lipgloss.NewStyle().Padding(0, 0).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4A32DE")).MarginLeft(1)
)
