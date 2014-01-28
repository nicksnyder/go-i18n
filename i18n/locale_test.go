package i18n

import (
	"testing"
)

func TestNewLocale(t *testing.T) {
	tests := []struct {
		localeId string
		valid    bool
	}{
		{"en-US", true},
		{"en_US", true},
		{"en US", false},
		{"en-US-en-US", false},
	}
	for _, test := range tests {
		_, err := NewLocale(test.localeId)
		if test.valid && err != nil {
			t.Errorf("%s should be a valid locale: %s", test.localeId, err)
		}
		if !test.valid && err == nil {
			t.Errorf("%s should not be a valid locale", test.localeId)
		}
	}
}
