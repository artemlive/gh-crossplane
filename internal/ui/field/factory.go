package field

import (
	"reflect"
	"strings"

	"github.com/artemlive/gh-crossplane/internal/util"
)

func GenerateComponentsByPaths(obj any, paths []string) []FieldComponent {
	root := reflect.ValueOf(obj)
	// unwrap pointer if necessary
	// it turns *GroupFile to GroupFile
	if root.Kind() == reflect.Ptr {
		root = root.Elem()
	}
	var components []FieldComponent

	for _, path := range paths {
		parts := strings.Split(path, ".") // e.g. "Spec.HasIssues" becomes ["Spec", "HasIssues"]
		v := root
		t := v.Type()

		// traverse each part (e.g. Spec -> HasIssues)
		for i, part := range parts {
			fieldVal := v.FieldByName(part)
			if !fieldVal.IsValid() {
				break // Skip invalid field
			}
			if i == len(parts)-1 {
				// we reached the last part, create a component
				structField, _ := t.FieldByName(part)
				tag := structField.Tag.Get("ui")
				meta := util.ParseTag(tag)

				switch meta["type"] {
				case "checkbox":
					if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Bool {
						components = append(components, NewCheckboxComponent(meta["label"], fieldVal.Interface().(*bool)))
					}
				case "text":
					if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.String {
						components = append(components, NewTextInputComponent(meta["label"], fieldVal.Interface().(*string)))
					} else if fieldVal.Kind() == reflect.String {
						components = append(components, NewTextInputComponent(meta["label"], fieldVal.Addr().Interface().(*string)))
					}
				}
			} else {
				// Go deeper
				if fieldVal.Kind() == reflect.Ptr {
					fieldVal = fieldVal.Elem()
				}
				v = fieldVal
				t = v.Type()
			}
		}
	}
	return components
}
