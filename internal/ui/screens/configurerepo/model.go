package configurerepo

import (
	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/ui/field"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
)

// compile-time check to ensure ConfigureRepoModel implements the ViewableModel interface
var _ ui.ViewableModel = (*ConfigureRepoModel)(nil)

type ConfigureRepoModel struct {
	repo         *domain.Repository
	fields       []field.FieldComponent
	focusedIndex int
	message      ui.Message
}

func New(repo *domain.Repository) *ConfigureRepoModel {
	comps := field.GenerateComponentsByPaths(repo, field.RepoEditableFields)
	if len(comps) > 0 {
		comps[0].Focus()
	}
	return &ConfigureRepoModel{
		repo:         repo,
		fields:       comps,
		focusedIndex: 0,
	}
}

func (m *ConfigureRepoModel) Init() tea.Cmd {
	return nil
}

func (m *ConfigureRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.String(); key {
		case "esc":
			return m, func() tea.Msg { return ui.SwitchToConfigureGroupMsg{} }

		case "up", "k":
			m.focusedIndex = (m.focusedIndex - 1 + len(m.fields)) % len(m.fields)
			m.fields[m.focusedIndex].Focus()
		case "down", "j":
			m.focusedIndex = (m.focusedIndex + 1) % len(m.fields)
			m.fields[m.focusedIndex].Focus()
		case "enter", "i":
			return m, func() tea.Msg {
				return ui.SwitchToGroupMsg{}
			}
		}
	case tea.WindowSizeMsg:
		ui.LastWindowSize = msg
	}

	updatedField, cmd := m.fields[m.focusedIndex].Update(msg, ui.ModeEditing)
	m.fields[m.focusedIndex] = updatedField
	return m, cmd
}

func (m *ConfigureRepoModel) View() (string, *tea.Cursor) {
	var out string
	var fields string
	for _, field := range m.fields {
		fields += field.View() + "\n"
	}
	if m.message.Msg != "" {
		out += "\n" + ui.FormatMessage(m.message)
	}

	return style.StyleModalBox(fields, ui.LastWindowSize.Width, ui.LastWindowSize.Height) + "\n" + out, nil
}
