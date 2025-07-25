package field

import (
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	tea "github.com/charmbracelet/bubbletea/v2"
)

// FieldComponent defines the interface for all field components
type FieldComponent interface {
	Init() tea.Cmd
	Update(tea.Msg, ui.FocusMode) (FieldComponent, tea.Cmd)
	View() string
	Focus() tea.Cmd
	Blur()
	IsFocused() bool
	Label() string
	// offset for the cursor position
	// e.g you have a label in the field and the cursor
	// should be placed after the label, so you need to return
	// the length of the label + 1 (for the space after the label)
	CursorOffset() int
}

// PreviewableComponent extends FieldComponent to support preview functionality
// The idea was to allow components to draw preview lines of their content
type PreviewableComponent interface {
	FieldComponent
	PreviewLines() []string
}

type Cursorer interface {
	Cursor() *tea.Cursor
}
type RenderResult struct {
	Components []RenderedComponent // main fields, for cursor tracking
	ExtraLines []string            // previews, modals, footers
}

type RenderedComponent struct {
	Component FieldComponent
	Lines     []string
}
