package configuregroup

import (
	"fmt"
	"reflect"
	"strings"
	"time"

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

// View renders the entire view of the ConfigureGroupModel.
// The cursor position must be calculated based on the curren layout
// the component itself knows it's X cursor offset, but the Y position is determined by the layout
func (m ConfigureGroupModel) View() (string, *tea.Cursor) {
	var layers []*lipgloss.Layer
	var globalCursor *tea.Cursor
	layoutY := 0

	// === TABS HEADER ===
	tabs := m.renderTabs()
	tabLines := strings.Split(tabs, "\n")
	layers = append(layers, lipgloss.NewLayer(tabs).Y(layoutY))
	layoutY += len(tabLines)

	// === SPACER ===
	layers = append(layers, lipgloss.NewLayer("").Y(layoutY))
	layoutY++

	// === DELEGATED RENDERING ===
	res := m.tabHandlers[m.activeTab].Render(&m)

	for _, rc := range res.Components {
		view := strings.Join(rc.Lines, "\n")
		layers = append(layers, lipgloss.NewLayer(view).Y(layoutY))

		// Track cursor if focused
		if rc.Component.IsFocused() {
			if c, ok := rc.Component.(field.Cursorer); ok {
				cursor := c.Cursor()
				if cursor != nil {
					globalCursor = tea.NewCursor(
						rc.Component.CursorOffset()+cursor.X,
						layoutY+cursor.Y,
					)
				}
			}
		}

		layoutY += len(rc.Lines)
	}

	// === EXTRAS (Previews, Modals etc.) ===
	for _, line := range res.ExtraLines {
		layers = append(layers, lipgloss.NewLayer(line).Y(layoutY))
		layoutY += lipgloss.Height(line)
	}

	// === STATUS BAR ===
	status := m.tabHandlers[m.activeTab].StatusBarText(&m)
	statusLines := strings.Split(status, "\n")
	layers = append(layers, lipgloss.NewLayer(status).Y(layoutY))
	layoutY += len(statusLines)

	// === MESSAGE ===
	if msg := m.renderMessage(); msg != "" {
		msgLines := strings.Split(msg, "\n")
		layers = append(layers, lipgloss.NewLayer(msg).Y(layoutY))
		layoutY += len(msgLines)
	}

	return lipgloss.NewCanvas(layers...).Render(), globalCursor
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
