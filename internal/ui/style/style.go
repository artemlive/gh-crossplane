package style

import (
	"strings"

	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/lipgloss/v2"
)

// TODO: consolidate this as a theme
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

	LabelStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("7")) // white
	InactiveTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))         // gray
	FocusedTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Bold(true) // bright green
	FieldBlockStyle   = lipgloss.NewStyle().Padding(0, 1)
	FocusedPrefix     = "➤"
	DimStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#555"))
)

var TextInputStyleEditingFocused = textinput.StyleState{
	Prompt: lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
	Text:   lipgloss.NewStyle().Foreground(lipgloss.Color("229")),
}

var TextInputStyleEditingBlurred = textinput.StyleState{
	Prompt: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	Text:   lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
}

var TextInputStyleNavigationFocused = textinput.StyleState{
	Prompt: lipgloss.NewStyle().Foreground(lipgloss.Color("81")),
	Text:   lipgloss.NewStyle().Foreground(lipgloss.Color("81")),
}

var TextInputStyleNavigationBlurred = textinput.StyleState{
	Prompt: lipgloss.NewStyle().Foreground(lipgloss.Color("238")),
	Text:   lipgloss.NewStyle().Foreground(lipgloss.Color("238")),
}

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

func DimBackground(lines []string, width int) string {
	var out []string
	for _, line := range lines {
		padded := line + strings.Repeat(" ", max(0, width-len(line)))
		out = append(out, DimStyle.Render(padded))
	}
	return strings.Join(out, "\n")
}

func RenderModalCentered(content string, width, height int) string {
	lines := strings.Split(content, "\n")
	contentHeight := len(lines)
	contentWidth := maxLineLength(lines)

	paddingTop := max(0, (height-contentHeight)/2)
	paddingLeft := max(0, (width-contentWidth)/2)

	var out []string
	for i := 0; i < paddingTop; i++ {
		out = append(out, "")
	}
	for _, line := range lines {
		out = append(out, strings.Repeat(" ", paddingLeft)+line)
	}
	return strings.Join(out, "\n")
}

func maxLineLength(lines []string) int {
	maxLen := 0
	for _, l := range lines {
		if len(l) > maxLen {
			maxLen = len(l)
		}
	}
	return maxLen
}
