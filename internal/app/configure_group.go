package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/artemlive/gh-crossplane/internal/manifest"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabHandler interface {
	Render(m *ConfigureGroupModel) []string
	Update(m *ConfigureGroupModel, msg tea.Msg) (tea.Model, tea.Cmd)
	StatusBarText(m *ConfigureGroupModel) string
}

type GenericTabHandler struct{}

func (h GenericTabHandler) Render(m *ConfigureGroupModel) []string {
	var lines []string

	components := m.fieldComponents[m.activeTab]
	for _, comp := range components {
		lines = append(lines, comp.View())
	}
	return lines

}

func (h GenericTabHandler) Update(m *ConfigureGroupModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case FieldDoneMsg:
		if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 {
			comps[m.focusedIndex].Blur()
		}
		m.mode = ModeNavigation
		return m, nil

	case FieldDoneUpMsg:
		return m.handlePrevField()

	case FieldDoneDownMsg:
		return m.handleNextField()

	case tea.KeyMsg:
		switch m.mode {
		case ModeNavigation:
			switch msg.String() {
			case "tab", "down", "j":
				return m.handleNextField()
			case "shift+tab", "up", "k":
				return m.handlePrevField()
			case "left", "h":
				return m.switchTab(-1)
			case "right", "l":
				return m.switchTab(1)
			case "enter", "i":
				if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 {
					comps[m.focusedIndex].Focus()
					m.mode = ModeEditing
				}
				return m, nil
			case "ctrl+s":
				if err := m.loader.SaveGroupFile(m.group); err != nil {
					m.message = ErrorMessage("Error saving group: " + err.Error())
				} else {
					m.message = InfoMessage(fmt.Sprintf("Group '%s' saved successfully.", m.group.Title()))
				}
				return m, nil
			case "esc":
				return m, func() tea.Msg { return switchToMenuMsg{} }
			case "q":
				return m, tea.Quit
			}

		case ModeEditing:
			if msg.String() == "esc" {
				return m.Update(FieldDoneMsg{})
			}
		}
	}

	// Pass input to focused field
	if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 && m.focusedIndex < len(comps) {
		newComp, cmd := comps[m.focusedIndex].Update(msg, m.mode)
		m.fieldComponents[m.activeTab][m.focusedIndex] = newComp
		return m, cmd
	}

	return m, nil
}

func (h GenericTabHandler) StatusBarText(m *ConfigureGroupModel) string {
	switch m.mode {
	case ModeNavigation:
		return configureGroupStatusStyleNavigation.Render("[NAV Mode] Up/Down Left/Right to navigate, Enter to edit, Ctrl+s to save, q to quit")
	case ModeEditing:
		return configureGroupStatusStyleEditing.Render("[EDT Mode] Press Esc or Enter to finish")
	}
	return ""
}

type RepositoryTabHandler struct{}

func (h RepositoryTabHandler) Render(m *ConfigureGroupModel) []string {
	var lines []string

	// Render the repository components
	components := m.fieldComponents[m.activeTab]
	for _, comp := range components {
		lines = append(lines, comp.View())

		if comp.IsFocused() {
			if repoComp, ok := comp.(*RepoComponent); ok {
				previewBlock := strings.Join(GenerateRepoPreviewLines(repoComp.repo), "\n")
				lines = append(lines, repoPreviewStyle.Render(previewBlock))
			}
		}
	}

	return lines
}

func (h *RepositoryTabHandler) Update(m *ConfigureGroupModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FieldDoneMsg:
		if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 {
			comps[m.focusedIndex].Blur()
		}
		m.mode = ModeNavigation
		return m, nil

	case FieldDoneUpMsg:
		return m.handlePrevField()

	case FieldDoneDownMsg:
		return m.handleNextField()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			return m.handlePrevField()
		case "down", "j":
			return m.handleNextField()
		case "left", "h":
			return m.switchTab(-1)
		case "right", "l":
			return m.switchTab(1)
		case "ctrl+s":
			err := m.loader.SaveGroupFile(m.group)
			if err != nil {
				m.message = ErrorMessage("Error saving group: " + err.Error())
			} else {
				m.message = InfoMessage(fmt.Sprintf("Group '%s' saved successfully.", m.group.Title()))
			}
			return m, nil
		case "esc":
			return m, func() tea.Msg { return switchToMenuMsg{} }
		}
	}

	// Forward input to focused component
	if comps := m.fieldComponents[m.activeTab]; len(comps) > 0 && m.focusedIndex < len(comps) {
		newComp, cmd := comps[m.focusedIndex].Update(msg, m.mode)
		m.fieldComponents[m.activeTab][m.focusedIndex] = newComp
		return m, cmd
	}

	return m, nil
}

func (h RepositoryTabHandler) StatusBarText(m *ConfigureGroupModel) string {
	return "[REPO Mode] Up/Down to navigate, Ctrl+s to save, q to quit"
}

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

func GenerateRepoComponents(repos []domain.Repository) []FieldComponent {
	var components []FieldComponent
	for i, repo := range repos {
		components = append(components, NewRepoComponent(repo, i))
	}
	return components
}

func GenerateRepoPreviewLines(obj any) []string {
	var lines []string

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Parse `ui` tag and skip if no label
		tag := field.Tag.Get("ui")
		meta := parseTag(tag)
		label := meta["label"]
		if label == "" {
			continue
		}

		// Skip complex types
		kind := value.Kind()
		if kind == reflect.Slice || kind == reflect.Map || kind == reflect.Struct {
			continue
		}

		// Unwrap pointer if needed
		if kind == reflect.Ptr {
			if value.IsNil() {
				// TODO: I think we should show only fields that are set
				// lines = append(lines, fmt.Sprintf("%s: <nil>", label))
				continue
			}
			value = value.Elem()
			kind = value.Kind()
		}

		// Format the value
		var str string
		switch kind {
		case reflect.String:
			// TODO: same here, we should show only if set
			if value.String() != "" {
				str = value.String()
			}
		case reflect.Bool:
			str = boolToStr(value.Bool())
		default:
			str = fmt.Sprintf("%v", value.Interface())
		}

		if str != "" {
			lines = append(lines, fmt.Sprintf("%s: %s", label, str))
		}
	}

	return lines
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

	tabHandlers []TabHandler
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

	// initialize field components for each tab
	m.tabHandlers = make([]TabHandler, len(FieldGroups))
	for i, fg := range FieldGroups {
		if fg.GroupLevel {
			components := GenerateComponentsByPaths(&group.Manifest, fg.FieldPaths)
			m.fieldComponents = append(m.fieldComponents, components)
			m.tabHandlers[i] = &GenericTabHandler{}
		} else {
			m.fieldComponents = append(m.fieldComponents, GenerateRepoComponents(group.Manifest.Spec.Repositories))
			m.tabHandlers[i] = &RepositoryTabHandler{}
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
	newTab := (m.activeTab + delta + len(m.tabs)) % len(m.tabs)
	m.activeTab = newTab
	if m.tabs[m.activeTab].TabName == "Repositories" {
		// force repositories tab to be in navigation mode
		m.mode = ModeNavigation
	}
	m.focusedIndex = 0
	return m, nil
}

func (m ConfigureGroupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.tabHandlers[m.activeTab].Update(&m, msg)
}

func (m ConfigureGroupModel) View() string {
	var lines []string

	lines = m.tabHandlers[m.activeTab].Render(&m)

	if len(lines) == 0 {
		lines = append(lines, "Not supported yet or no fields available in this tab.")
	}

	// Status bar rendering
	statusBar := m.tabHandlers[m.activeTab].StatusBarText(&m)

	// Optional message display (info, warning, error)
	message := m.renderMessage()

	return m.renderTabs() + "\n\n" + strings.Join(lines, "\n\n") + "\n\n" + statusBar + "\n" + message
}

func (m ConfigureGroupModel) renderMessage() string {
	if m.message.Msg == "" {
		return ""
	}

	var styledMsg string
	switch m.message.Type {
	case MessageTypeInfo:
		styledMsg = infoMessageStyle.Render("[Info] " + m.message.Msg)
	case MessageTypeError:
		styledMsg = errorMessageStyle.Render("[Error] " + m.message.Msg)
	case MessageTypeWarning:
		styledMsg = warningMessageStyle.Render("[Warning] " + m.message.Msg)
	}
	return "\n\n" + styledMsg
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
