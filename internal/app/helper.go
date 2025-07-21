package app

import (
	"reflect"
	"strings"
)

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

func ptrToBool(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return v.Elem().Bool()
	}
	return false
}

func boolToStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func ifEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}
