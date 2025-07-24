package field

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/debug"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type TextInputComponent struct {
	label   string
	value   *string
	ti      textinput.Model
	focused bool
	mode    ui.FocusMode
}

func NewTextInputComponent(label string, val *string) *TextInputComponent {
	ti := textinput.New()
	ti.VirtualCursor = false
	// to fit the placeholder
	ti.SetWidth(20)

	if val != nil {
		ti.SetValue(*val)
	}

	ti.Styles = textinput.DefaultStyles(true)

	return &TextInputComponent{
		label:   label,
		value:   val,
		ti:      ti,
		focused: false,
	}
}

func (c *TextInputComponent) SetValue(val string) {
	if c.value == nil {
		c.value = new(string)
	}
	*c.value = val
	c.ti.SetValue(val)
}

func (c *TextInputComponent) Value() string {
	if c.value == nil {
		return ""
	}
	return *c.value
}

func (c *TextInputComponent) SetPlaceholder(placeholder string) {
	c.ti.Placeholder = placeholder
}

func (c *TextInputComponent) Cursor() *tea.Cursor {
	if c.ti.VirtualCursor {
		return nil
	}
	return c.ti.Cursor()
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
			c.ti.Styles.Focused = style.TextInputStyleEditingFocused
		} else {
			c.ti.Styles.Blurred = style.TextInputStyleEditingBlurred
		}

	case ui.ModeNavigation:
		if c.focused {
			c.ti.Styles.Focused = style.TextInputStyleNavigationFocused
		} else {
			c.ti.Styles.Blurred = style.TextInputStyleNavigationBlurred
		}

	default:
		c.ti.Styles.Blurred = style.TextInputStyleEditingBlurred
	}
}

func (c *TextInputComponent) SetLabel(label string) {
	c.label = label
}

func (c *TextInputComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	c.mode = mode
	c.applyStyles()

	if mode == ui.ModeEditing {
		var cmd tea.Cmd
		c.ti, cmd = c.ti.Update(msg)
		c.SetValue(c.ti.Value())

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

	// navigation mode: ignore typing but allow blinking
	if _, ok := msg.(tea.KeyMsg); ok {
		return c, nil
	}

	var cmd tea.Cmd
	c.ti, cmd = c.ti.Update(msg)
	return c, cmd
}

func (c *TextInputComponent) CursorOffset() int {
	return lipgloss.Width(c.label) + 2 // +2 for ": "
}

func (c *TextInputComponent) Focus() tea.Cmd {
	debug.Log.Printf("Focus() called on field '%s'", c.label)
	c.focused = true
	return c.ti.Focus()
}

func (c *TextInputComponent) Blur() {
	debug.Log.Printf("Blur() called on field '%s'", c.label)
	c.focused = false
	c.ti.Blur()
}

func (c *TextInputComponent) IsFocused() bool { return c.focused }

func (c *TextInputComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (c *TextInputComponent) Label() string {
	return c.label
}
