package style

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

	ModalBoxStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Background(lipgloss.Color("234")).
			Align(lipgloss.Center)

	TextInputStyleEditingFocused    = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
	TextInputStyleEditingBlurred    = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	TextInputStyleNavigationFocused = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	TextInputStyleNavigationBlurred = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
)

//	func StyleModalBox(content string, width, height int) string {
//		box := ModalBoxStyle.Render(content)
//
//		// Center the modal in the available space
//		return lipgloss.Place(width, height,
//			lipgloss.Center, lipgloss.Center,
//			box,
//		)
//	}
func StyleModalBox(content string, termWidth, termHeight int) string {
	const modalWidth = 60
	const modalHeight = 15

	lines := strings.Split(content, "\n")

	// pad lines to modalWidth manually
	for i, line := range lines {
		lines[i] = lipgloss.NewStyle().
			Width(modalWidth).
			Render(line)
	}

	// pad vertical if not enough lines
	for len(lines) < modalHeight {
		lines = append(lines, strings.Repeat(" ", modalWidth))
	}

	contentBlock := lipgloss.JoinVertical(lipgloss.Top, lines[:modalHeight]...)

	box := ModalBoxStyle.
		Width(modalWidth).
		Height(modalHeight).
		Render(contentBlock)

	return lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, box)
}
