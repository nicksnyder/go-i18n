package bundle

import (
	"fmt"
	"reflect"
)

// Whether the Value represents a struct value.
func isStruct(v reflect.Value) bool {
	return v.Kind() == reflect.Struct
}

// Converts a struct to a map[string]interface{}.
// Pulled from github.com/fatih/structs, modified to omit tagging and
// nested conversion.
func structToMap(v reflect.Value) (map[string]interface{}, error) {
	if !isStruct(v) {
		return nil, fmt.Errorf("struct expected, got %s", v.Kind())
	}

	out := make(map[string]interface{})
	fields, err := structFieldNames(v)
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		out[field] = v.FieldByName(field).Interface()
	}
	return out, nil
}

// Creates a slice of the struct's field names.
// Pulled from github.com/fatih/structs, modified to omit tagging and to
// only return the field name rather than the full reflect.StructField.
func structFieldNames(v reflect.Value) ([]string, error) {
	if !isStruct(v) {
		return nil, fmt.Errorf("struct expected, got %s", v.Kind())
	}

	var fields []string
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// unexported field. skip.
		if field.PkgPath != "" {
			continue
		}
		fields = append(fields, field.Name)
	}
	return fields, nil
}
