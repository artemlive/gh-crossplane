package field

import (
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	tea "github.com/charmbracelet/bubbletea"
)

// FieldComponent defines the interface for all field components
type FieldComponent interface {
	Init() tea.Cmd
	Update(tea.Msg, ui.FocusMode) (FieldComponent, tea.Cmd)
	View() string
	Focus()
	Blur()
	IsFocused() bool
}
