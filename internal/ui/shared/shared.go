package shared

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
