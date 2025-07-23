package field

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/debug"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	view := c.ti.View()
	debug.Log.Printf("Rendering [%s]: %s, focused: %t", c.label, view, c.focused)
	return fmt.Sprintf("%s: %s", c.label, view)
}

func (c *TextInputComponent) applyStyles() {
	switch c.mode {
	case ui.ModeEditing:
		if c.focused {
			c.ti.PromptStyle = style.TextInputStyleEditingFocused
			c.ti.Cursor.Style = style.TextInputStyleEditingFocused
		} else {
			c.ti.PromptStyle = style.TextInputStyleEditingBlurred
			c.ti.Cursor.Style = style.TextInputStyleEditingBlurred
		}
	case ui.ModeNavigation:
		if c.focused {
			c.ti.PromptStyle = style.TextInputStyleNavigationFocused
			c.ti.Cursor.Style = style.TextInputStyleNavigationFocused
		} else {
			c.ti.PromptStyle = style.TextInputStyleNavigationBlurred
			c.ti.Cursor.Style = style.TextInputStyleNavigationBlurred
		}
	default:
		c.ti.PromptStyle = style.TextInputStyleEditingBlurred
		c.ti.Cursor.Style = style.TextInputStyleEditingBlurred
	}
}
func (c *TextInputComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	c.mode = mode
	c.applyStyles()

	// Handle input depending on mode
	if mode == ui.ModeEditing {
		var cmd tea.Cmd
		c.ti, cmd = c.ti.Update(msg)

		if c.value != nil {
			*c.value = c.ti.Value()
		}

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
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

	// Navigation mode: ignore typing but allow blinking
	if _, ok := msg.(tea.KeyMsg); ok {
		return c, nil
	}

	var cmd tea.Cmd
	c.ti, cmd = c.ti.Update(msg)
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
