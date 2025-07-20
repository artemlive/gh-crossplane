package app

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
