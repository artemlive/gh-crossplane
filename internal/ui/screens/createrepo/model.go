package createrepo

import (
	"fmt"

	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateRepoModel struct {
	step        int
	repoName    string
	description string

	message string
	input   textinput.Model
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

func (m CreateRepoModel) Init() tea.Cmd {
	return nil
}

func (m CreateRepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := m.input.Value()
			if val == "" {
				m.message = "please enter a value"
				return m, nil // Do nothing if input is empty
			}

			switch m.step {
			case StepRepoName:
				m.repoName = val
				m.step = StepDescription
				m.input.SetValue("")
			case StepDescription:
				m.description = val
				m.step = StepDone
				m.input.SetValue("")
				return m, func() tea.Msg {
					return ui.SwitchToSelectGroupMsg{
						RepoName:    m.repoName,
						Description: m.description,
					}
				}
			}
		default:
			m.message = ""
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m CreateRepoModel) View() string {
	var prompt string
	switch m.step {
	case StepRepoName:
		prompt = "Enter repository name:"
		m.input.Placeholder = "repo-name"
	case StepDescription:
		prompt = "Enter repository description:"
		m.input.Placeholder = "description"
	case StepDone:
		prompt = "Done!"
	}

	information := ""
	if m.message != "" {
		information = fmt.Sprintf("\n\n%s", style.InfoMessageStyle.Render(m.message))
	}

	return fmt.Sprintf("%s\n%s\n%s", prompt, m.input.View(), information)
}
