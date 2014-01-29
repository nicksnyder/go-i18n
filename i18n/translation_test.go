package i18n

import (
	"testing"
)

func TestNewSingleTranslation(t *testing.T) {
	t.Skipf("not implemented")
}

func testNewTranslation(t *testing.T, data map[string]interface{}) Translation {
	translation, err := NewTranslation(data)
	if err != nil {
		t.Fatal(err)
	}
	return translation
}
