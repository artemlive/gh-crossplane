package app

import (
	"fmt"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/manifest"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	ModeMenu mode = iota
	ModeCreateRepo
	ModeSelectGroup
	ModeDone
)

type model struct {
	currentMode mode

	menu            MenuModel
	createRepoForm  CreateRepoModel
	selectGroupForm SelectGroupModel

	createdRepo   domain.Repository
	selectedGroup string
	msg           string

	manifestLoader *manifest.ManifestLoader
}

func NewAppModel(groupDir string) model {
	return model{
		currentMode:    ModeMenu,
		menu:           NewMenuModel(),
		manifestLoader: manifest.NewManifestLoader(groupDir),
	}
}

func (m model) Init() tea.Cmd {
	return m.menu.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentMode {
	case ModeMenu:
		newMenu, cmd, choice := m.menu.Update(msg)
		m.menu = newMenu
		switch choice {
		case ChoiceCreateRepo:
			m.createRepoForm = NewCreateRepoModel()
			m.currentMode = ModeCreateRepo
		}
		return m, cmd

	case ModeCreateRepo:
		newForm, cmd, done := m.createRepoForm.Update(msg)
		m.createRepoForm = newForm
		if done {
			groupNames := m.manifestLoader.Groups()
			if len(groupNames) == 0 {
				m.msg = fmt.Sprintf("warning: no groups found")
				m.currentMode = ModeDone
				return m, nil
			}
			m.selectGroupForm = NewSelectGroupModel(groupNames)
			m.currentMode = ModeSelectGroup
		}
		return m, cmd

	case ModeSelectGroup:
		newSelect, cmd, selected, done := m.selectGroupForm.Update(msg)
		m.selectGroupForm = newSelect
		if done {
			m.createdRepo = domain.Repository{
				Name:        m.createRepoForm.repoName,
				Description: m.createRepoForm.description,
			}
			m.selectedGroup = selected
			// Inject the createdRepo into the group and write to file
			gf := m.manifestLoader.GetGroup(selected)
			if gf == nil {
				m.msg = fmt.Sprintf("failed to load group %s", selected)
				m.currentMode = ModeDone
				return m, nil
			}
			gf.Manifest.Spec.Repositories = append(gf.Manifest.Spec.Repositories, m.createdRepo)
			if err := m.manifestLoader.SaveGroupFile(gf); err != nil {
				m.msg = fmt.Sprintf("failed to save group file: %v", err)
			} else {
				m.msg = fmt.Sprintf("Repo %q added to group %q", m.createdRepo.Name, selected)
			}
			m.currentMode = ModeDone
		}
		return m, cmd

	case ModeDone:
		if key, ok := msg.(tea.KeyMsg); ok && (key.String() == "enter" || key.String() == "q") {
			m.currentMode = ModeMenu
		}
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	switch m.currentMode {
	case ModeMenu:
		return m.menu.View()
	case ModeCreateRepo:
		return m.createRepoForm.View()
	case ModeSelectGroup:
		return m.selectGroupForm.View()
	case ModeDone:
		return fmt.Sprintf("\n%s\n\nPress enter or q to continue...", m.msg)
	default:
		return "unknown mode"
	}
}
