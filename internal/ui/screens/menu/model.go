package menu

import (
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	tea "github.com/charmbracelet/bubbletea"
)

type MenuModel struct {
	choices map[MenuChoice]string
	cursor  int
}

type MenuChoice int

const (
	ChoiceCreateRepo MenuChoice = iota
	ChoiceConfigureRepo
)

var menuLabels = map[MenuChoice]string{
	ChoiceCreateRepo:    "Create a new repo",
	ChoiceConfigureRepo: "Configure an existing repo",
}

var menuOrder = []MenuChoice{
	ChoiceCreateRepo,
	ChoiceConfigureRepo,
}

func NewMenuModel() MenuModel {
	return MenuModel{
		choices: menuLabels,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch menuOrder[m.cursor] {
			case ChoiceCreateRepo:
				return m, func() tea.Msg { return ui.SwitchToCreateRepoMsg{} }
			case ChoiceConfigureRepo:
				return m, func() tea.Msg { return ui.SwitchToSelectGroupMsg{} }
			}
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := "What do you want to do?\n\n"
	for i, choice := range menuOrder {
		cursor := " "
		if m.cursor == i {
			cursor = style.MainMenuCurosrStyle.Render(">")
			s += cursor + " " + style.MainMenuCurosrStyle.Render(menuLabels[choice]) + "\n"
		} else {
			s += cursor + " " + menuLabels[choice] + "\n"
		}
	}
	s += "\nUse up/down to move, enter to select, q to quit."
	return s
}
