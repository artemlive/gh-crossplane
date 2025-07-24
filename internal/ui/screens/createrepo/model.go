package createrepo

import (
	"strings"

	"github.com/artemlive/gh-crossplane/debug"
	"github.com/artemlive/gh-crossplane/internal/ui/field"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type CreateRepoModel struct {
	step        int
	repoName    string
	description string

	message string
	input   *field.TextInputComponent
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
	ti := field.NewTextInputComponent("Repository Name", nil)
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
			debug.Log.Printf("Input %v, val: %s", m.input, val)
			if val == "" {
				m.message = "please enter a value"
				return m, nil // Do nothing if input is empty
			}

			switch m.step {
			case StepRepoName:
				m.repoName = val
				m.step = StepDescription
				m.input.SetValue("") // Clear input for next step
			case StepDescription:
				m.description = val
				m.step = StepDone
				return m, func() tea.Msg {
					return ui.SwitchToSelectGroupMsg{
						RepoName:    m.repoName,
						Description: m.description,
					}
				}
			}
		case "esc":
			//TODO: switch between steps
			return m, func() tea.Msg {
				return ui.SwitchToMenuMsg{}
			}
		default:
			m.message = ""
		}
	}
	newInput, cmd := m.input.Update(msg, ui.ModeEditing)
	m.input = newInput.(*field.TextInputComponent)

	return m, cmd
}

func (m CreateRepoModel) View() (string, *tea.Cursor) {
	var layers []*lipgloss.Layer
	var globalCursor *tea.Cursor
	layoutY := 0

	// Prompt
	var prompt string
	switch m.step {
	case StepRepoName:
		m.input.SetPlaceholder("repo-name")
	case StepDescription:
		m.input.SetLabel("Description")
		m.input.SetPlaceholder("description")
	case StepDone:
		prompt = "Done!"
	}

	m.input.Focus()
	promptLines := strings.Split(prompt, "\n")
	layers = append(layers, lipgloss.NewLayer(prompt).Y(layoutY))
	layoutY += len(promptLines)

	// Input field
	inputView := m.input.View()
	inputLines := strings.Split(inputView, "\n")
	layers = append(layers, lipgloss.NewLayer(inputView).Y(layoutY))

	// Compute global cursor
	if cur := m.input.Cursor(); cur != nil {
		globalCursor = tea.NewCursor(m.input.CursorOffset()+cur.X, layoutY+cur.Y)
	}
	layoutY += len(inputLines)

	// info message (if any)
	if m.message != "" {
		msg := style.InfoMessageStyle.Render(m.message)
		msgLines := strings.Split(msg, "\n")
		layers = append(layers, lipgloss.NewLayer("\n\n"+msg).Y(layoutY))
		layoutY += 2 + len(msgLines) //nolint
	}

	// Final canvas
	canvas := lipgloss.NewCanvas(layers...)
	return canvas.Render(), globalCursor
}
