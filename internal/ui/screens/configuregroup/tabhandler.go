package configuregroup

import (
	"fmt"
	"strings"
	"time"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/ui/field"
	"github.com/artemlive/gh-crossplane/internal/ui/screens/configurerepo"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	tea "github.com/charmbracelet/bubbletea"
)

type TabHandler interface {
	Render(m *ConfigureGroupModel) []string
	Update(m *ConfigureGroupModel, msg tea.Msg) (tea.Model, tea.Cmd)
	StatusBarText(m *ConfigureGroupModel) string
}

type GenericTabHandler struct{}

func (h GenericTabHandler) Render(m *ConfigureGroupModel) []string {
	var lines []string

	components := m.fieldComponents[m.activeTab]
	for _, comp := range components {
		lines = append(lines, comp.View())
	}
	return lines

}

func (h GenericTabHandler) Update(m *ConfigureGroupModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case field.FieldDoneMsg:
		m.mode = ui.ModeNavigation
		// force re-rendering the style if needed
		return m, func() tea.Msg { return ui.TickMsg{} }

	case ui.TickMsg:
		// update all fields for blinking
		for i, comp := range m.fieldComponents[m.activeTab] {
			newComp, _ := comp.Update(msg, m.mode)
			m.fieldComponents[m.activeTab][i] = newComp
		}

		// set up the next tick
		// we need this to update the fields for blinking
		nextTick := tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return ui.TickMsg{}
		})
		return m, nextTick

	case field.FieldDoneUpMsg:
		return m.handlePrevField()

	case field.FieldDoneDownMsg:
		return m.handleNextField()

	case tea.KeyMsg:
		switch m.mode {
		case ui.ModeNavigation:
			switch msg.String() {
			case "tab", "down", "j":
				return m.handleNextField()
			case "shift+tab", "up", "k":
				return m.handlePrevField()
			case "left", "h":
				return m.switchTab(-1)
			case "right", "l":
				return m.switchTab(1)
			case "enter", "i":
				var cmd tea.Cmd
				if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 {
					// blur the old focused field
					comps[m.focusedIndex].Blur()

					// focus the current one
					cmd = comps[m.focusedIndex].Focus()
					m.mode = ui.ModeEditing
				}
				return m, cmd
			case "ctrl+s":
				if err := m.loader.SaveGroupFile(m.group); err != nil {
					m.message = ui.ErrorMessage(fmt.Sprintf("Error saving group '%s': %s", m.group.Title(), err.Error()))
				} else {
					m.message = ui.InfoMessage(fmt.Sprintf("Group '%s' saved successfully.", m.group.Title()))
				}
				return m, nil
			case "esc":
				return m, func() tea.Msg { return ui.SwitchToMenuMsg{} }
			case "q":
				return m, tea.Quit
			}

		case ui.ModeEditing:
			if msg.String() == "esc" {
				return m.Update(field.FieldDoneMsg{})
			}
		}
	}
	var cmd tea.Cmd
	for i := range m.fieldComponents[m.activeTab] {
		if i == m.focusedIndex {
			newComp, c := m.fieldComponents[m.activeTab][i].Update(msg, m.mode)
			m.fieldComponents[m.activeTab][i] = newComp
			cmd = c
		}
	}
	return m, cmd
}

func (h GenericTabHandler) StatusBarText(m *ConfigureGroupModel) string {
	switch m.mode {
	case ui.ModeNavigation:
		return style.ConfigureGroupStatusStyleNavigation.Render("[NAV Mode] Up/Down Left/Right to navigate, Enter to edit, Ctrl+s to save, q to quit")
	case ui.ModeEditing:
		return style.ConfigureGroupStatusStyleEditing.Render("[EDT Mode] Press Esc or Enter to finish")
	}
	return ""
}

type RepositoryTabHandler struct{}

func (h RepositoryTabHandler) Render(m *ConfigureGroupModel) []string {
	if m.repoModal != nil {
		// Modal takes over render it exclusively
		return []string{m.repoModal.View()}
	}

	var lines []string

	// Render the repository components
	components := m.fieldComponents[m.activeTab]
	for _, comp := range components {
		lines = append(lines, comp.View())

		if comp.IsFocused() {
			if pv, ok := comp.(field.PreviewableComponent); ok {
				lines = append(lines, style.RepoPreviewStyle.Render(strings.Join(pv.PreviewLines(), "\n")))
			}
		}
	}

	return lines
}

func (h *RepositoryTabHandler) Update(m *ConfigureGroupModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isModalOpen() {
		newModel, cmd := m.repoModal.Update(msg)
		m.repoModal = newModel
		return m, cmd
	}

	switch msg := msg.(type) {

	case ui.SwitchToGroupMsg:
		m.repoModal = nil
		return m, nil

	case field.FieldDoneMsg:
		if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 {
			comps[m.focusedIndex].Blur()
		}
		m.mode = ui.ModeNavigation
		return m, nil

	case field.FieldDoneUpMsg:
		return m.handlePrevField()

	case field.FieldDoneDownMsg:
		return m.handleNextField()
	case field.FieldOpenMsg:
		switch v := msg.Value.(type) {
		case *domain.Repository:
			m.repoModal = configurerepo.New(v)
		default:
			m.message = ui.ErrorMessage("Invalid repository value type in FieldOpenMsg")
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			err := m.loader.SaveGroupFile(m.group)
			if err != nil {
				m.message = ui.ErrorMessage("Error saving group: " + err.Error())
			} else {
				m.message = ui.InfoMessage(fmt.Sprintf("Group '%s' saved successfully.", m.group.Title()))
			}
			return m, nil
		case "esc":
			return m, func() tea.Msg { return ui.SwitchToMenuMsg{} }

		case "q":
			return m, tea.Quit
		}
	}

	// Let component try to handle it first
	if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 && m.focusedIndex < len(comps) {
		newComp, cmd := comps[m.focusedIndex].Update(msg, m.mode)
		m.fieldComponents[m.activeTab][m.focusedIndex] = newComp

		if cmd != nil {
			return m, cmd
		}
	}

	// Only handle navigation *after* component had a chance
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "up", "k":
			return m.handlePrevField()
		case "down", "j":
			return m.handleNextField()
		case "left", "h":
			return m.switchTab(-1)
		case "right", "l":
			return m.switchTab(1)
		}
	}

	return m, nil
}
func (h RepositoryTabHandler) StatusBarText(m *ConfigureGroupModel) string {
	return "[REPO Mode] Up/Down to navigate, Ctrl+s to save, q to quit"
}
