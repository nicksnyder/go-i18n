package bundle

import (
	"reflect"
	"testing"
)

type testPerson struct {
	Person string
	Age    int
}

var (
	fieldNames = []string{"Person", "Age"}
	testMap    = map[string]interface{}{
		"Person": "Bob",
		"Age":    26,
	}
	testStruct = testPerson{Person: "Bob", Age: 26}
	testValue  = reflect.ValueOf(testStruct)
)

func TestStructToMap(t *testing.T) {
	testMatchesMap(t, "with struct", testMap, structToMap(testValue))
}

func TestStructToMapPanicsWithNonStruct(t *testing.T) {
	defer testPanics(t, "struct expected, got ptr")
	testMatchesMap(t, "with pointer", testMap, structToMap(reflect.ValueOf(&testStruct)))
}

func TestStructFieldNames(t *testing.T) {
	fields := structFieldNames(testValue)
	if !reflect.DeepEqual(fields, fieldNames) {
		t.Errorf("expected %s, got %s", fieldNames, fields)
	}
}

func TestStructFieldNamesPanicsWithNonStruct(t *testing.T) {
	defer testPanics(t, "struct expected, got ptr")
	structFieldNames(reflect.ValueOf(&testStruct))
}

func testPanics(t *testing.T, expected string) {
	if e := recover(); e != nil {
		if err, ok := e.(error); ok {
			if err.Error() != expected {
				t.Errorf("expected '%s', got '%s'",
					expected,
					err.Error())
			}
		} else {
			t.Errorf("expected error, got: %v", e)
		}
	} else {
		t.Errorf("expected error, none found")
	}
}
