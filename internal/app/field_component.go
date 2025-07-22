package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CheckboxComponent struct {
	Label   string
	Value   *bool
	Focused bool
}

func NewCheckboxComponent(label string, ptr *bool) *CheckboxComponent {
	return &CheckboxComponent{
		Label: label,
		Value: ptr,
	}
}

func (c *CheckboxComponent) Init() tea.Cmd {
	return nil
}

func (c *CheckboxComponent) Update(msg tea.Msg, mode FocusMode) (FieldComponent, tea.Cmd) {
	if mode != ModeEditing {
		return c, nil // only handle updates in editing mode
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "enter":
			if c.Value == nil {
				c.Value = new(bool)
				// init to true if it was nil
				*c.Value = true
			} else {
				*c.Value = !*c.Value
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
	if c.Value != nil && *c.Value {
		checked = "x"
	}
	cursor := " "
	if c.Focused {
		cursor = ">"
	}
	return fmt.Sprintf("%s [%s] %s", cursor, checked, c.Label)
}

func (c *CheckboxComponent) Focus() {
	c.Focused = true
}

func (c *CheckboxComponent) Blur() {
	c.Focused = false
}

func (c *CheckboxComponent) IsFocused() bool {
	return c.Focused
}

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

func (c *TextInputComponent) Update(msg tea.Msg, mode FocusMode) (FieldComponent, tea.Cmd) {
	if mode != ModeEditing {
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

// FieldComponent defines the interface for all field components
type FieldComponent interface {
	Init() tea.Cmd
	Update(tea.Msg, FocusMode) (FieldComponent, tea.Cmd)
	View() string
	Focus()
	Blur()
	IsFocused() bool
}

type FieldDoneMsg struct{}
type FieldDoneUpMsg struct{}
type FieldDoneDownMsg struct{}

func GenerateComponentsByPaths(obj any, paths []string) []FieldComponent {
	root := reflect.ValueOf(obj)
	// unwrap pointer if necessary
	// it turns *GroupFile to GroupFile
	if root.Kind() == reflect.Ptr {
		root = root.Elem()
	}
	var components []FieldComponent

	for _, path := range paths {
		parts := strings.Split(path, ".") // e.g. "Spec.HasIssues" becomes ["Spec", "HasIssues"]
		v := root
		t := v.Type()

		// traverse each part (e.g. Spec -> HasIssues)
		for i, part := range parts {
			fieldVal := v.FieldByName(part)
			if !fieldVal.IsValid() {
				break // Skip invalid field
			}
			if i == len(parts)-1 {
				// we reached the last part, create a component
				structField, _ := t.FieldByName(part)
				tag := structField.Tag.Get("ui")
				meta := parseTag(tag)

				switch meta["type"] {
				case "checkbox":
					if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Bool {
						components = append(components, NewCheckboxComponent(meta["label"], fieldVal.Interface().(*bool)))
					}
				case "text":
					if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.String {
						components = append(components, NewTextInputComponent(meta["label"], fieldVal.Interface().(*string)))
					} else if fieldVal.Kind() == reflect.String {
						components = append(components, NewTextInputComponent(meta["label"], fieldVal.Addr().Interface().(*string)))
					}
				}
			} else {
				// Go deeper
				if fieldVal.Kind() == reflect.Ptr {
					fieldVal = fieldVal.Elem()
				}
				v = fieldVal
				t = v.Type()
			}
		}
	}
	return components
}

// parseTag parses a struct tag string into a map of key-value pairs.
// e.g. "label=Has Issues, type=checkbox" becomes
// map[string]string{"label": "Has Issues", "type": "checkbox"}
func parseTag(tag string) map[string]string {
	parts := strings.Split(tag, ",")
	out := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return out
}

type RepoComponent struct {
	repo    *domain.Repository
	index   int
	focused bool
}

func NewRepoComponent(repo domain.Repository, index int) *RepoComponent {
	return &RepoComponent{
		repo:  &repo,
		index: index,
	}
}

func (r *RepoComponent) View() string {
	prefix := "  "
	if r.focused {
		prefix = "→ "
	}
	name := ifEmpty(r.repo.Name, "<unnamed>")
	return prefix + name
}

func (r *RepoComponent) Focus() { r.focused = true }
func (r *RepoComponent) Blur()  { r.focused = false }
func (r *RepoComponent) Update(msg tea.Msg, mode FocusMode) (FieldComponent, tea.Cmd) {
	return r, nil
}

func (r *RepoComponent) IsFocused() bool {
	return r.focused
}

func (r *RepoComponent) Init() tea.Cmd {
	return nil
}
