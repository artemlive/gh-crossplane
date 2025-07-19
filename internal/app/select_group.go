package app

import (
	"os"

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
	os.Exit(1)
	return nil
}

func NewSelectGroupModel(groups []manifest.GroupFile) SelectGroupModel {
	listKeys := newListKeyMap()

	listItems := make([]list.Item, len(groups))
	for i, group := range groups {
		listItems[i] = group
	}

	// set the defaul size for the list due to bug if I set it to 0,0
	groupsList := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
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
		onStartup:  true,
	}
}

func (m SelectGroupModel) Update(msg tea.Msg) (SelectGroupModel, tea.Cmd, string, bool) {
	//	_, isWindowSizeMsg := msg.(tea.WindowSizeMsg)
	//
	//	// Since this program is using the full size of the viewport we
	//	// need to wait until we've received the window dimensions before
	//	// we can initialize the viewport. The initial dimensions come in
	//	// quickly, though asynchronously, which is why we wait for them
	//	// here.
	//	if m.onStartup && !isWindowSizeMsg {
	//		return m, nil, "", false
	//	}
	//	if m.onStartup && isWindowSizeMsg {
	//		h, v := appStyle.GetFrameSize()
	//		m.list.SetSize(msg.(tea.WindowSizeMsg).Width-h, msg.(tea.WindowSizeMsg).Height-v)
	//		m.onStartup = false
	//	}
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
				return m, nil, "", false
			}
			m.selectedGroup = m.groupNames[m.cursor].Manifest.Metadata.Name
		}

	}

	newSelectGroupModel, cmd := m.list.Update(msg)
	m.list = newSelectGroupModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...), m.selectedGroup, false

}

func (m SelectGroupModel) View() string {
	return appStyle.Render(m.list.View())
}
