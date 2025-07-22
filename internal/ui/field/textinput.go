package field

import (
	"fmt"

	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputComponent struct {
	label   string
	value   *string
	ti      textinput.Model
	focused bool
}

func NewTextInputComponent(label string, ptr *string) *TextInputComponent {
	ti := textinput.New()
	if ptr != nil {
		ti.SetValue(*ptr)
	}
	ti.Focus()
	return &TextInputComponent{label: label, value: ptr, ti: ti, focused: true}
}

func (c *TextInputComponent) View() string {
	return fmt.Sprintf("%s: %s", c.label, c.ti.View())
}

func (c *TextInputComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	if mode != ui.ModeEditing {
		return c, nil // Only handle updates in editing mode
	}
	var cmd tea.Cmd
	c.ti, cmd = c.ti.Update(msg)

	val := c.ti.Value()
	if c.value != nil {
		*c.value = val
	} else {
		c.value = new(string)
		*c.value = val
	}

	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		return c, func() tea.Msg { return FieldDoneMsg{} }
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

func (c *TextInputComponent) Focus()          { c.ti.Focus(); c.focused = true }
func (c *TextInputComponent) Blur()           { c.ti.Blur(); c.focused = false }
func (c *TextInputComponent) IsFocused() bool { return c.focused }
func (c *TextInputComponent) Init() tea.Cmd {
	return nil
}
