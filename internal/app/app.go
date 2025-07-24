package app

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/debug"
	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/manifest"
	"github.com/artemlive/gh-crossplane/internal/ui/screens/configuregroup"
	"github.com/artemlive/gh-crossplane/internal/ui/screens/createrepo"
	"github.com/artemlive/gh-crossplane/internal/ui/screens/menu"
	"github.com/artemlive/gh-crossplane/internal/ui/screens/selectgroup"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	tea "github.com/charmbracelet/bubbletea/v2"
)

type appState struct {
	manifestLoader *manifest.ManifestLoader
}

func (m *appState) GetManifestLoader() *manifest.ManifestLoader {
	return m.manifestLoader
}

type model struct {
	curScreen ui.ViewableModel
	state     appState

	message ui.Message
	width   int
	height  int
}

func NewAppModel(groupDir string) model {
	state := appState{
		manifestLoader: manifest.NewManifestLoader(groupDir),
	}

	return model{
		state:     state,
		curScreen: menu.NewMenuModel(),
	}
}

func switchToMenu() tea.Cmd {
	return func() tea.Msg {
		return ui.SwitchToMenuMsg{}
	}
}

func (m model) Init() tea.Cmd {
	return switchToMenu()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ui.SwitchToMenuMsg:
		menu := menu.NewMenuModel()
		m.curScreen = menu
		return m, menu.Init()
	case ui.SwitchToCreateRepoMsg:
		createRepoModel := createrepo.NewCreateRepoModel()
		m.curScreen = createRepoModel
		return m, createRepoModel.Init()
	case ui.SwitchToSelectGroupMsg:
		// pass the width and height to the selectGroup model
		// because it needs to know the size of the terminal
		// since the window size message wasn't sent there on initialization
		repo := domain.Repository{
			Name:        msg.RepoName,
			Description: msg.Description,
		}
		selectGroupModel := selectgroup.NewSelectGroupModel(m.state.GetManifestLoader().Groups(), repo, m.width, m.height)
		m.curScreen = selectGroupModel
		return m, selectGroupModel.Init()
	case ui.SwitchToConfigureGroupMsg:
		groupName := msg.GroupName
		group := m.state.GetManifestLoader().GetGroup(groupName)
		if group == nil {
			// TODO: render error message on the screen
			m.message = ui.ErrorMessage(fmt.Sprintf("Group '%s' not found", groupName))
			return m, nil
		}
		configureGroupModel := configuregroup.NewConfigureGroupModel(group, m.state.manifestLoader, m.width, m.height)
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
	debug.Log.Printf("Delegating message: %T", msg)
	// Delegate message to current screen
	model, cmd := m.curScreen.Update(msg)
	screen, ok := model.(ui.ViewableModel)
	if !ok {
		m.message = ui.ErrorMessage("Current screen does not implement ViewableModel")
	}
	m.curScreen = screen
	return m, cmd
}

func (m model) View() (string, *tea.Cursor) {
	if m.curScreen == nil {
		debug.Log.Println("No current screen to render")
		return "No screen to display", nil
	}

	debug.Log.Printf("Rendering screen: %T", m.curScreen)
	return m.curScreen.View()
}
