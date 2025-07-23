package field

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/debug"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
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
	ti.VirtualCursor = false

	if ptr != nil {
		ti.SetValue(*ptr)
	}

	ti.Styles = textinput.DefaultStyles(true)

	return &TextInputComponent{
		label:   label,
		value:   ptr,
		ti:      ti,
		focused: false,
	}
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
	debug.Log.Printf("Applying styles for field '%s' in mode '%v', cursor blinking: %t", c.label, c.mode, c.ti.Cursor().Blink)
	c.ti.Cursor().Color = lipgloss.Color("red")
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

func (c *TextInputComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	c.mode = mode
	c.applyStyles()

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

	// navigation mode: ignore typing but allow blinking
	if _, ok := msg.(tea.KeyMsg); ok {
		return c, nil
	}

	var cmd tea.Cmd
	c.ti, cmd = c.ti.Update(msg)
	return c, cmd
}

func (c *TextInputComponent) Focus() tea.Cmd {
	debug.Log.Printf("Focus() called on field '%s'", c.label)
	c.focused = true
	c.ti.Cursor().Blink = true
	return c.ti.Focus()
}

func (c *TextInputComponent) Blur() {
	debug.Log.Printf("Blur() called on field '%s'", c.label)
	c.focused = false
	c.ti.Blur()
}

func (c *TextInputComponent) IsFocused() bool { return c.focused }

func (c *TextInputComponent) Init() tea.Cmd {
	return nil
}
