package configuregroup

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/artemlive/gh-crossplane/debug"
	"github.com/artemlive/gh-crossplane/internal/manifest"
	"github.com/artemlive/gh-crossplane/internal/ui/field"
	ui "github.com/artemlive/gh-crossplane/internal/ui/shared"
	"github.com/artemlive/gh-crossplane/internal/ui/style"
	"github.com/artemlive/gh-crossplane/internal/util"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type MessageType int

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
		meta := util.ParseTag(tag)
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
			str = util.BoolToStr(value.Bool())
		default:
			str = fmt.Sprintf("%v", value.Interface())
		}

		if str != "" {
			lines = append(lines, fmt.Sprintf("%s: %s", label, str))
		}
	}

	return lines
}

func (m *ConfigureGroupModel) isModalOpen() bool {
	return m.modal != nil
}

type ConfigureGroupModel struct {
	tabs            []field.FieldGroup
	fieldComponents [][]field.FieldComponent // one slice per tab
	activeTab       int
	group           *manifest.GroupFile
	repoIndex       int // which repo is selected in "Repositories" tab
	width           int
	height          int
	focusedIndex    int

	mode    ui.FocusMode // current focus mode, either navigation or editing
	loader  *manifest.ManifestLoader
	message ui.Message

	tabHandlers []TabHandler
	modal       ui.ViewableModel
}

func NewConfigureGroupModel(group *manifest.GroupFile, loader *manifest.ManifestLoader, width, height int) *ConfigureGroupModel {
	m := ConfigureGroupModel{
		tabs:         field.FieldGroups,
		activeTab:    0,
		group:        group,
		repoIndex:    0,
		width:        width,
		height:       height,
		focusedIndex: 0,
		loader:       loader,
	}

	// initialize field components for each tab
	m.tabHandlers = make([]TabHandler, len(field.FieldGroups))
	for i, fg := range field.FieldGroups {
		if fg.GroupLevel {
			components := field.GenerateComponentsByPaths(&group.Manifest, fg.FieldPaths)
			m.fieldComponents = append(m.fieldComponents, components)
			m.tabHandlers[i] = &GenericTabHandler{}
		} else {
			repoComponent := field.NewRepositoriesComponent("Repositories", &group.Manifest.Spec.Repositories)
			m.fieldComponents = append(m.fieldComponents, []field.FieldComponent{repoComponent})
			m.tabHandlers[i] = &RepositoryTabHandler{}
		}
	}
	return &m
}

func (m *ConfigureGroupModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	// focus first field if present
	comps := m.fieldComponents[m.activeTab]
	if len(comps) > 0 {
		m.focusedIndex = 0
		cmds = append(cmds, comps[0].Focus())
	}

	// tick used to update fields for blinking
	tick := tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return ui.TickMsg{}
	})
	cmds = append(cmds, tick)

	return tea.Batch(cmds...)
}

func (m *ConfigureGroupModel) handleNextField() (tea.Model, tea.Cmd) {
	components := m.fieldComponents[m.activeTab]
	if len(components) == 0 {
		return m, nil
	}
	components[m.focusedIndex].Blur()
	m.focusedIndex = (m.focusedIndex + 1) % len(components)
	return m, components[m.focusedIndex].Focus()
}

func (m *ConfigureGroupModel) handlePrevField() (tea.Model, tea.Cmd) {
	components := m.fieldComponents[m.activeTab]
	if len(components) == 0 {
		return m, nil
	}
	components[m.focusedIndex].Blur()
	m.focusedIndex = (m.focusedIndex - 1 + len(components)) % len(components)

	return m, components[m.focusedIndex].Focus()
}

func (m *ConfigureGroupModel) switchTab(delta int) (tea.Model, tea.Cmd) {
	m.activeTab = (m.activeTab + delta + len(m.tabs)) % len(m.tabs)

	if m.tabs[m.activeTab].TabName == "Repositories" {
		// force repositories tab to be in navigation mode
		m.mode = ui.ModeNavigation
	}

	comps := m.fieldComponents[m.activeTab]
	m.focusedIndex = 0
	var cmd tea.Cmd
	if len(comps) > 0 {
		cmd = comps[0].Focus()
	}

	return m, cmd
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

	statusBar := m.tabHandlers[m.activeTab].StatusBarText(&m)
	message := m.renderMessage()

	return m.renderTabs() + "\n\n" + strings.Join(lines, "\n\n") + "\n\n" + statusBar + "\n" + message
}

func (m ConfigureGroupModel) Cursor() *tea.Cursor {
	debug.Log.Printf("we are in cursor method\n")
	var comps = m.fieldComponents[m.activeTab]
	if m.focusedIndex < len(comps) {
		if c, ok := comps[m.focusedIndex].(field.Cursorer); ok && comps[m.focusedIndex].IsFocused() {
			cursor := c.Cursor()
			debug.Log.Printf("cursor = %+v", cursor)
			return cursor
		}
	}
	return nil
}

func (m ConfigureGroupModel) renderMessage() string {
	if m.message.Msg == "" {
		return ""
	}

	return "\n\n" + ui.FormatMessage(m.message)
}

func (m ConfigureGroupModel) renderTabs() string {
	var rendered []string
	for i, tab := range m.tabs {
		curStyle := style.InactiveTabStyle
		if i == m.activeTab {
			curStyle = style.ActiveTabStyle
		}
		rendered = append(rendered, curStyle.Render(tab.TabName))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}
