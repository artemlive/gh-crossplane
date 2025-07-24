package field

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/internal/domain"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	"github.com/artemlive/gh-crossplane/internal/util"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// RepositoriesComponent is a list of repositories with preview and open/edit support
type RepositoriesComponent struct {
	label   string
	repos   *[]domain.Repository
	index   int // currently focused index
	focused bool
}

func NewRepositoriesComponent(label string, repos *[]domain.Repository) *RepositoriesComponent {
	return &RepositoriesComponent{
		label:   label,
		repos:   repos,
		index:   0,
		focused: false,
	}
}

func (c *RepositoriesComponent) View() string {
	lines := []string{style.LabelStyle.Render(c.label + ":")}

	if len(*c.repos) == 0 {
		lines = append(lines, style.InactiveTextStyle.Render("No repositories"))
		return style.FieldBlockStyle.Render(ui.JoinVertical(lines))
	}

	for i, repo := range *c.repos {
		prefix := "  "
		if c.focused && c.index == i {
			prefix = fmt.Sprintf("%s ", style.FocusedPrefix)
			lines = append(lines, style.FocusedTextStyle.Render(prefix+repo.Name))
		} else {
			lines = append(lines, prefix+repo.Name)
		}
	}

	return style.FieldBlockStyle.Render(ui.JoinVertical(lines))
}

func (c *RepositoriesComponent) Update(msg tea.Msg, mode ui.FocusMode) (FieldComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if c.index > 0 {
				c.index--
			}
		case "down", "j":
			if c.index < len(*c.repos)-1 {
				c.index++
			}
		case "enter":
			return c, func() tea.Msg {
				return FieldOpenMsg{
					Value: &(*c.repos)[c.index],
					Label: c.label,
				}

			}
		}
	}
	return c, nil
}

func (c *RepositoriesComponent) Focus() tea.Cmd {
	c.focused = true
	return nil
}

func (c *RepositoriesComponent) Blur() {
	c.focused = false
}

func (c *RepositoriesComponent) IsFocused() bool {
	return c.focused
}

func (c *RepositoriesComponent) Init() tea.Cmd {
	return nil
}

func (c *RepositoriesComponent) Label() string {
	return c.label
}

func (c *RepositoriesComponent) CursorOffset() int {
	if len(*c.repos) == 0 || c.index >= len(*c.repos) {
		return 0
	}
	return lipgloss.Width((*c.repos)[c.index].Name) + 2 // +2 for the prefix and space
}

func (c *RepositoriesComponent) PreviewLines() []string {
	if len(*c.repos) == 0 || c.index >= len(*c.repos) {
		return nil
	}

	r := (*c.repos)[c.index]
	var lines []string

	if r.Name != "" {
		lines = append(lines, "Name: "+r.Name)
	}
	if r.Description != "" {
		lines = append(lines, "Description: "+r.Description)
	}
	if r.Visibility != "" {
		lines = append(lines, "Visibility: "+r.Visibility)
	}
	if r.DefaultBranch != "" {
		lines = append(lines, "Default Branch: "+r.DefaultBranch)
	}
	if r.Archived != nil {
		lines = append(lines, "Archived: "+util.BoolToStr(*r.Archived))
	}
	if r.AllowAutoMerge != nil {
		lines = append(lines, "Allow Auto-Merge: "+util.BoolToStr(*r.AllowAutoMerge))
	}
	if r.DeleteBranchOnMerge != nil {
		lines = append(lines, "Delete Branch on Merge: "+util.BoolToStr(*r.DeleteBranchOnMerge))
	}

	return lines
}
