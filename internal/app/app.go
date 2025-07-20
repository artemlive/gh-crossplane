package app

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/manifest"
	tea "github.com/charmbracelet/bubbletea"
)

type switchToMenuMsg struct{}
type switchToCreateRepoMsg struct{}
type switchToSelectGroupMsg struct {
	repoName    string
	description string
}

type switchToConfigureGroupMsg struct {
	// maybe we could return the group object instead of just the name?
	groupName string
}

type appState struct {
	manifestLoader *manifest.ManifestLoader
}

func (m *appState) GetManifestLoader() *manifest.ManifestLoader {
	return m.manifestLoader
}

type model struct {
	curScreen tea.Model
	state     appState

	message Message
	width   int
	height  int
}

func NewAppModel(groupDir string) model {
	state := appState{
		manifestLoader: manifest.NewManifestLoader(groupDir),
	}

	return model{
		state:     state,
		curScreen: NewMenuModel(),
	}
}

func switchToMenu() tea.Cmd {
	return func() tea.Msg {
		return switchToMenuMsg{}
	}
}

func (m model) Init() tea.Cmd {
	return switchToMenu()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case switchToMenuMsg:
		menu := NewMenuModel()
		m.curScreen = menu
		return m, menu.Init()
	case switchToCreateRepoMsg:
		createRepoModel := NewCreateRepoModel()
		m.curScreen = createRepoModel
		return m, createRepoModel.Init()
	case switchToSelectGroupMsg:
		// pass the width and height to the selectGroup model
		// because it needs to know the size of the terminal
		// since the window size message wasn't sent there on initialization
		repo := domain.Repository{
			Name:        msg.repoName,
			Description: msg.description,
		}
		selectGroupModel := NewSelectGroupModel(m.state.GetManifestLoader().Groups(), repo, m.width, m.height)
		m.curScreen = selectGroupModel
		return m, selectGroupModel.Init()
	case switchToConfigureGroupMsg:
		groupName := msg.groupName
		group := m.state.GetManifestLoader().GetGroup(groupName)
		if group == nil {
			// TODO: render error message on the screen
			m.message = ErrorMessage(fmt.Sprintf("Group '%s' not found", groupName))
			return m, nil
		}
		configureGroupModel := NewConfigureGroupModel(group, m.state.manifestLoader, m.width, m.height)
		m.curScreen = configureGroupModel
		return m, configureGroupModel.Init()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// I save this to pass the size to the selectGroup model
		m.width = msg.Width
		m.height = msg.Height

	}
	// Delegate message to current screen
	newScreen, cmd := m.curScreen.Update(msg)
	m.curScreen = newScreen
	return m, cmd
}

func (m model) View() string {
	if m.curScreen == nil {
		return "No screen to display"
	}

	return m.curScreen.View()
}
