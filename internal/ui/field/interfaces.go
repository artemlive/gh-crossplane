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
	Focus() tea.Cmd
	Blur()
	IsFocused() bool
}

// PreviewableComponent extends FieldComponent to support preview functionality
// The idea was to allow components to draw preview lines of their content
type PreviewableComponent interface {
	FieldComponent
	PreviewLines() []string
}
