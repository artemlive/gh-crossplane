package app

import (
	"github.com/artemlive/gh-crossplane/internal/manifest"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectGroupModel struct {
	groupNames    []manifest.GroupFile
	cursor        int
	selectedGroup string
	list          list.Model
	keys          *listKeyMap

	onStartup bool
}

type listKeyMap struct {
	addGroup    key.Binding
	selectGroup key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		addGroup: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add new group"),
		),
		selectGroup: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select group"),
		),
	}
}

func (m SelectGroupModel) Init() tea.Cmd {
	return nil
}

func NewSelectGroupModel(groups []manifest.GroupFile, width, height int) SelectGroupModel {
	listKeys := newListKeyMap()

	listItems := make([]list.Item, len(groups))
	for i, group := range groups {
		listItems[i] = group
	}

	// set the defaul size for the list due to bug if I set it to 0,0
	groupsList := list.New(listItems, list.NewDefaultDelegate(), width, height)
	groupsList.Title = "Select Group Or Add New by pressing 'a'"
	groupsList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.addGroup,
			listKeys.selectGroup,
		}
	}
	return SelectGroupModel{
		groupNames: groups,
		cursor:     0,
		list:       groupsList,
		keys:       listKeys,
	}
}

func (m SelectGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.selectGroup):
			if len(m.groupNames) == 0 {
				return m, nil
			}
			m.selectedGroup = m.groupNames[m.cursor].Manifest.Metadata.Name
		}

	}

	newSelectGroupModel, cmd := m.list.Update(msg)
	m.list = newSelectGroupModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

}

func (m SelectGroupModel) View() string {
	return appStyle.Render(m.list.View())
}
