package util

import (
	"reflect"
	"strings"
)

func PtrToBool(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return v.Elem().Bool()
	}
	return false
}

func BoolToStr(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func IfEmpty(s, fallback string) string {
	if strings.TrimSpace(s) == "" {
		return fallback
	}
	return s
}

// ParseTag parses a struct tag string into a map of key-value pairs.
// e.g. "label=Has Issues, type=checkbox" becomes
// map[string]string{"label": "Has Issues", "type": "checkbox"}
func ParseTag(tag string) map[string]string {
	parts := strings.Split(tag, ",")
	out := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return out
}
