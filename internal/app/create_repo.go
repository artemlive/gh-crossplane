package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateRepoModel struct {
	step        int
	repoName    string
	description string
	curTeamName string
	permission  string
	teamPerms   []TeamPermission

	input textinput.Model
}

type TeamPermission struct {
	Team       string
	Permission string
}

const (
	StepRepoName = iota
	StepDescription
	StepTeamName
	StepPermission
	StepAskAddMore
	StepDone
)

func NewCreateRepoModel() CreateRepoModel {
	ti := textinput.New()
	ti.Placeholder = "repo-name"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40
	return CreateRepoModel{
		step:  StepRepoName,
		input: ti,
	}
}

func (m CreateRepoModel) Update(msg tea.Msg) (CreateRepoModel, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := m.input.Value()
			switch m.step {
			case StepRepoName:
				m.repoName = val
				m.step = StepDescription
				m.input.SetValue("")
				m.input.Placeholder = "repo-name"
			case StepDescription:
				m.description = val
				m.step = StepDone
				m.input.SetValue("")
				m.input.Placeholder = "description"
				// we are done, so we can proceed to group selection
				return m, nil, true
			}
		case "esc":
			return m, tea.Quit, false
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd, false
}

func (m CreateRepoModel) View() string {
	var prompt string
	switch m.step {
	case StepRepoName:
		prompt = "Enter repository name:"
	case StepDescription:
		prompt = "Enter repository description:"
	case StepDone:
		prompt = "Done!"
	}

	summary := ""
	for _, tp := range m.teamPerms {
		summary += fmt.Sprintf(" - %s: %s\n", tp.Team, tp.Permission)
	}

	return fmt.Sprintf("%s\n%s\n\n%s", prompt, m.input.View(), summary)
}
