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
	actual, err := structToMap(testValue)
	if err != nil {
		t.Errorf("not expecting error, got %v", err)
	}
	testMatchesMap(t, "with struct", testMap, actual)
}

func TestStructToMapErrorsWithNonStruct(t *testing.T) {
	_, err := structToMap(reflect.ValueOf(&testStruct))
	if err == nil {
		t.Errorf("expecting expecting error, got nothing", err)
	}
}

func TestStructFieldNames(t *testing.T) {
	fields, err := structFieldNames(testValue)
	if err != nil {
		t.Errorf("not expecting error, got %v", err)
	}
	if !reflect.DeepEqual(fields, fieldNames) {
		t.Errorf("expected %s, got %s", fieldNames, fields)
	}
}

func TestStructFieldNamesErrorsWithNonStruct(t *testing.T) {
	_, err := structFieldNames(reflect.ValueOf(&testStruct))
	if err == nil {
		t.Errorf("expecting expecting error, got nothing", err)
	}
}
