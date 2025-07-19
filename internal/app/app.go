package app

import (
	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/manifest"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

type screenDoneMsg struct{}
type switchToMenuMsg struct{}
type switchToCreateRepoMsg struct{}
type switchToSelectGroupMsg struct {
	repoName    string
	description string
}

type appState struct {
	manifestLoader *manifest.ManifestLoader
	createdRepo    domain.Repository
	selectedGroup  string
	message        string
}

func (m *appState) GetManifestLoader() *manifest.ManifestLoader {
	return m.manifestLoader
}

type model struct {
	curScreen tea.Model
	state     appState

	width  int
	height int
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
		selectGroupModel := NewSelectGroupModel(m.state.GetManifestLoader().Groups(), m.width, m.height)
		m.curScreen = selectGroupModel
		return m, selectGroupModel.Init()
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
