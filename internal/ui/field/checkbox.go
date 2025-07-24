package field

import (
	"fmt"

	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	tea "github.com/charmbracelet/bubbletea/v2"
)

type CheckboxComponent struct {
	label   string
	value   *bool
	Focused bool
}

func NewCheckboxComponent(label string, ptr *bool) *CheckboxComponent {
	return &CheckboxComponent{
		label: label,
		value: ptr,
	}
}

func (c *CheckboxComponent) Init() tea.Cmd {
	return nil
}

func (c *CheckboxComponent) Label() string {
	return c.label
}

func (c *CheckboxComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	if mode != ui.ModeEditing {
		return c, nil // only handle updates in editing mode
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "space", "enter":
			if c.value == nil {
				c.value = new(bool)
				// init to true if it was nil
				*c.value = true
			} else {
				*c.value = !*c.value
			}
		case "up", "k":
			// Move focus to the previous component
			return c, func() tea.Msg { return FieldDoneUpMsg{} }
		case "down", "j":
			// Move focus to the next component
			return c, func() tea.Msg { return FieldDoneDownMsg{} }
		}
	}
	return c, nil
}

func (c *CheckboxComponent) View() string {
	checked := " "
	if c.value != nil && *c.value {
		checked = "x"
	}
	cursor := " "
	if c.Focused {
		cursor = style.FocusedPrefix
	}
	return fmt.Sprintf("%s [%s] %s", cursor, checked, c.label)
}

func (c *CheckboxComponent) Focus() tea.Cmd {
	c.Focused = true
	return nil
}

func (c *CheckboxComponent) CursorOffset() int {
	return 0 // No cursor offset needed for checkbox
}

func (c *CheckboxComponent) Blur() {
	c.Focused = false
}

func (c *CheckboxComponent) IsFocused() bool {
	return c.Focused
}
