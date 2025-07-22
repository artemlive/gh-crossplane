package field

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/debug"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextInputComponent struct {
	label   string
	value   *string
	ti      textinput.Model
	focused bool
	mode    ui.FocusMode
}

func NewTextInputComponent(label string, ptr *string) *TextInputComponent {
	ti := textinput.New()
	if ptr != nil {
		ti.SetValue(*ptr)
	}
	return &TextInputComponent{label: label, value: ptr, ti: ti, focused: false}
}

func (c *TextInputComponent) View() string {
	// Set style based on mode + focus
	switch c.mode {
	case ui.ModeEditing:
		c.ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	case ui.ModeNavigation:
		if c.focused {
			c.ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // focused but not editable
		} else {
			c.ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // unfocused
		}
	}

	view := c.ti.View()
	debug.Log.Printf("Rendering [%s]: %s, focused: %t", c.label, view, c.focused)
	return fmt.Sprintf("%s: %s", c.label, view)
}

func (c *TextInputComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	c.mode = mode
	if mode != ui.ModeEditing {
		if _, ok := msg.(tea.KeyMsg); ok {
			return c, nil
		}
		var cmd tea.Cmd
		c.ti, cmd = c.ti.Update(msg)
		return c, cmd
	}

	var cmd tea.Cmd
	c.ti, cmd = c.ti.Update(msg)

	// Sync value if necessary
	if c.value != nil {
		*c.value = c.ti.Value()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return c, func() tea.Msg { return FieldDoneMsg{} }
		case "up":
			return c, func() tea.Msg { return FieldDoneUpMsg{} }
		case "down":
			return c, func() tea.Msg { return FieldDoneDownMsg{} }
		}
	}
	return c, cmd
}
func (c *TextInputComponent) Focus() tea.Cmd {
	debug.Log.Printf("Focus() called on field '%s'", c.label)

	cmd := c.ti.Focus()
	c.focused = true
	c.ti.Cursor.Blink = true
	return cmd
}

func (c *TextInputComponent) Blur() {
	debug.Log.Printf("Blur() called on field '%s'", c.label)

	c.ti.Blur()
	c.focused = false
}

func (c *TextInputComponent) IsFocused() bool { return c.focused }
func (c *TextInputComponent) Init() tea.Cmd {
	return nil
}
