package app

import (
	"strings"

	"github.com/artemlive/gh-crossplane/internal/manifest"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FieldGroup struct {
	TabName     string
	FieldPaths  []string // Dot-separated, e.g., "Spec.Visibility"
	GroupLevel  bool     // true = group, false = repo
	Description string   // optional
}

type FocusMode int

const (
	ModeNavigation FocusMode = iota // navigating between fields
	ModeEditing                     // editing a field
)

type MessageType int

const (
	MessageTypeInfo MessageType = iota
	MessageTypeError
	MessageTypeWarning
)

type Message struct {
	Msg  string
	Type MessageType
}

// TODO: consider moving this to a separate package or file
// this is a list of field groups that will be used to render the configuration UI
var FieldGroups = []FieldGroup{
	{
		TabName: "Group: General",
		FieldPaths: []string{
			"Spec.Visibility",
			"Spec.DefaultBranch",
			"Spec.Topics",
			"Spec.ArchiveOnDestroy",
			"Spec.AutoInit",
			"Spec.IsTemplate",
			"Spec.ManagementPolicies",
			"Spec.DeletionPolicy",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Features",
		FieldPaths: []string{
			"Spec.HasIssues",
			"Spec.HasDownloads",
			"Spec.HasWiki",
			"Spec.HasDiscussions",
			"Spec.AllowAutoMerge",
			"Spec.AllowSquashMerge",
			"Spec.AllowMergeCommit",
			"Spec.AllowRebaseMerge",
			"Spec.AllowUpdateBranch",
			"Spec.DeleteBranchOnMerge",
			"Spec.VulnerabilityAlerts",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Merge Messages",
		FieldPaths: []string{
			"Spec.MergeCommitMessage",
			"Spec.MergeCommitTitle",
			"Spec.SquashMergeCommitMessage",
			"Spec.SquashMergeCommitTitle",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Protections",
		FieldPaths: []string{
			"Spec.Protections",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Security",
		FieldPaths: []string{
			"Spec.SecurityAndAnalysis",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Autolinks",
		FieldPaths: []string{
			"Spec.AutolinkReferences",
		},
		GroupLevel: true,
	},
	{
		TabName: "Group: Permissions",
		FieldPaths: []string{
			"Spec.Permissions",
		},
		GroupLevel: true,
	},
	{
		TabName: "Repositories",
		FieldPaths: []string{
			"Spec.Repositories",
		},
		GroupLevel: false,
	},
}

type ConfigureGroupModel struct {
	tabs            []FieldGroup
	fieldComponents [][]FieldComponent // one slice per tab
	activeTab       int
	group           *manifest.GroupFile
	repoIndex       int // which repo is selected in "Repositories" tab
	width           int
	height          int
	focusedIndex    int

	mode    FocusMode // current focus mode, either navigation or editing
	loader  *manifest.ManifestLoader
	message Message
}

func NewConfigureGroupModel(group *manifest.GroupFile, loader *manifest.ManifestLoader, width, height int) ConfigureGroupModel {
	m := ConfigureGroupModel{
		tabs:         FieldGroups,
		activeTab:    0,
		group:        group,
		repoIndex:    0,
		width:        width,
		height:       height,
		focusedIndex: 0,
		loader:       loader,
	}

	// Initialize field components for each tab
	for _, fg := range FieldGroups {
		if fg.GroupLevel {
			components := GenerateComponentsByPaths(&group.Manifest, fg.FieldPaths)
			m.fieldComponents = append(m.fieldComponents, components)
		} else {
			m.fieldComponents = append(m.fieldComponents, nil) // placeholder for "Repositories"
		}
	}
	return m
}

func (m ConfigureGroupModel) Init() tea.Cmd {
	return nil
}

func (m *ConfigureGroupModel) handleNextField() (tea.Model, tea.Cmd) {
	components := m.fieldComponents[m.activeTab]
	if len(components) == 0 {
		return m, nil
	}
	components[m.focusedIndex].Blur()
	m.focusedIndex = (m.focusedIndex + 1) % len(components)
	components[m.focusedIndex].Focus()
	return m, nil
}

func (m *ConfigureGroupModel) handlePrevField() (tea.Model, tea.Cmd) {
	components := m.fieldComponents[m.activeTab]
	if len(components) == 0 {
		return m, nil
	}
	components[m.focusedIndex].Blur()
	m.focusedIndex = (m.focusedIndex - 1 + len(components)) % len(components)
	components[m.focusedIndex].Focus()
	return m, nil
}

func (m *ConfigureGroupModel) switchTab(delta int) (tea.Model, tea.Cmd) {
	newTab := m.activeTab + delta
	if newTab >= 0 && newTab < len(m.tabs) {
		m.activeTab = newTab
		m.focusedIndex = 0
	}
	return m, nil
}

func (m ConfigureGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case FieldDoneMsg:
		// Leave editing mode
		if m.activeTab < len(m.fieldComponents) {
			comps := m.fieldComponents[m.activeTab]
			if len(comps) > 0 {
				comps[m.focusedIndex].Blur()
			}
		}
		m.mode = ModeNavigation
		return m, nil

	case tea.KeyMsg:
		switch m.mode {

		case ModeNavigation:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit

			case "tab", "down", "j":
				return m.handleNextField()

			case "shift+tab", "up", "k":
				return m.handlePrevField()

			case "right", "l":
				return m.switchTab(1)

			case "left", "h":
				return m.switchTab(-1)

			case "enter", "i":
				if m.activeTab < len(m.fieldComponents) {
					comps := m.fieldComponents[m.activeTab]
					if len(comps) > 0 {
						comp := comps[m.focusedIndex]
						comp.Focus()
						m.mode = ModeEditing
					}
				}
				return m, nil
			case "ctrl+s":
				// Save the group configuration
				if err := m.loader.SaveGroupFile(m.group); err != nil {
					m.message = ErrorMessage("Error saving group: " + err.Error())
				} else {
					m.message = InfoMessage("Group saved successfully.")
				}
				return m, nil
			}

		case ModeEditing:
			// Escape from editing
			if msg.String() == "esc" {
				return m.Update(FieldDoneMsg{})
			}
		}
	}

	// Forward input to focused component
	if m.activeTab < len(m.fieldComponents) {
		comps := m.fieldComponents[m.activeTab]
		if len(comps) > 0 && m.focusedIndex < len(comps) {
			newComp, cmd := comps[m.focusedIndex].Update(msg, m.mode)
			m.fieldComponents[m.activeTab][m.focusedIndex] = newComp
			return m, cmd
		}
	}

	return m, nil
}

func (m ConfigureGroupModel) View() string {
	var lines []string

	// show field components if this tab has them
	if m.activeTab < len(m.fieldComponents) && m.fieldComponents[m.activeTab] != nil {
		for _, c := range m.fieldComponents[m.activeTab] {
			lines = append(lines, c.View())
		}
	} else if m.tabs[m.activeTab].TabName == "Repositories" { // we assume that the Repositories tab is a special case
		for i, r := range m.group.Manifest.Spec.Repositories {
			selector := "  "
			if i == m.repoIndex {
				selector = "→ "
			}
			name := ifEmpty(r.Name, "<unnamed>")
			lines = append(lines, selector+name)
		}
	}
	var statusBar strings.Builder
	style := lipgloss.NewStyle()
	switch m.mode {
	case ModeNavigation:
		statusBar.WriteString("[Navigation Mode] Up/Down to navigate, Enter to edit, Left/Right to switch tabs, q to quit")
		style = configureGroupStatusStyleNavigation
	case ModeEditing:
		statusBar.WriteString("[Editing Mode] Press Esc or Enter to finish")
		style = configureGroupStatusStyleEditing
	}
	if len(lines) == 0 {
		lines = append(lines, "Not supported yet or no fields available in this tab.")
	}
	renderedView := m.renderTabs() + "\n\n" + strings.Join(lines, "\n\n") + "\n\n" + style.Render(statusBar.String())

	message := ""
	if m.message.Msg != "" {
		style := lipgloss.NewStyle()
		switch m.message.Type {
		case MessageTypeInfo:
			style = infoMessageStyle
			message = style.Render("[Info] " + m.message.Msg)
		case MessageTypeError:
			style = errorMessageStyle
			message = style.Render("[Error] " + m.message.Msg)
		case MessageTypeWarning:
			style = warningMessageStyle
			message = style.Render("[Warning] " + m.message.Msg)
		}
		renderedView += "\n\n" + message
	}
	return renderedView

}

func (m ConfigureGroupModel) renderTabs() string {
	var rendered []string
	for i, tab := range m.tabs {
		style := inactiveTabStyle
		if i == m.activeTab {
			style = activeTabStyle
		}
		rendered = append(rendered, style.Render(tab.TabName))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func boolStr(b *bool) string {
	if b == nil {
		return "unset"
	}
	if *b {
		return "true"
	}
	return "false"
}

func ifEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

var (
	inactiveTabStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(lipgloss.Color("240"))
	activeTabStyle   = inactiveTabStyle.BorderBottom(false).Bold(true)
)

// TODO: add switchToConfigureRepoMsg for per-repo editing
