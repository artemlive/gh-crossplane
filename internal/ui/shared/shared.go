package shared

import (
	"github.com/artemlive/gh-crossplane/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

type FocusMode int

const (
	ModeNavigation FocusMode = iota // navigating between fields
	ModeEditing                     // editing a field
)

type MessageType int

const (
	MessageTypeError   MessageType = iota // error message
	MessageTypeWarning                    // warning message
	MessageTypeInfo                       // informational message
)

type Message struct {
	Msg  string
	Type MessageType
}

var LastWindowSize struct {
	Width  int
	Height int
}

func ErrorMessage(msg string) Message {
	return Message{
		Msg:  msg,
		Type: MessageTypeError,
	}
}

func WarningMessage(msg string) Message {
	return Message{
		Msg:  msg,
		Type: MessageTypeWarning,
	}
}

func InfoMessage(msg string) Message {
	return Message{
		Msg:  msg,
		Type: MessageTypeInfo,
	}
}

type SwitchToMenuMsg struct{}
type SwitchToCreateRepoMsg struct{}
type SwitchToSelectGroupMsg struct {
	RepoName    string
	Description string
}

type SwitchToConfigureGroupMsg struct {
	GroupName string
}

type SwitchToGroupMsg struct {
	Index int
	Repo  *domain.Repository
}

type ForceReRenderMsg struct{}
type TickMsg struct{}

func FormatMessage(msg Message) string {
	switch msg.Type {
	case MessageTypeError:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("✖ " + msg.Msg)
	case MessageTypeWarning:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("⚠ " + msg.Msg)
	case MessageTypeInfo:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render("ℹ " + msg.Msg)
	default:
		return msg.Msg
	}
}
