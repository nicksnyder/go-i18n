package bundle

import (
	"fmt"
	"reflect"
)

// Panic if the value does not represent a struct.
func assertStruct(v reflect.Value) {
	if kind := v.Kind(); kind != reflect.Struct {
		panic(fmt.Errorf("struct expected, got %s", kind))
	}
}

// Converts a struct to a map[string]interface{}.
// Pulled from github.com/fatih/structs, modified to omit tagging and
// nested conversion.
func structToMap(v reflect.Value) map[string]interface{} {
	assertStruct(v)

	out := make(map[string]interface{})
	fields := structFieldNames(v)
	for _, field := range fields {
		out[field] = v.FieldByName(field).Interface()
	}
	return out
}

// Creates a slice of the struct's field names.
// Pulled from github.com/fatih/structs, modified to omit tagging and to
// only return the field name rather than the full reflect.StructField.
func structFieldNames(v reflect.Value) []string {
	assertStruct(v)

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
	return fields
}
