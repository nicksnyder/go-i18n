package locale

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		localeID string
		valid    bool
	}{
		{"en-US", true},
		{"en_US", true},
		{"zh-Hans-CN", true},
		{"zh-Hant-TW", true},
		{"zh-CN", false},
		{"zh-TW", false},
		{"en US", false},
		{"en-US-en-US", false},
	}
	for _, test := range tests {
		_, err := New(test.localeID)
		if test.valid && err != nil {
			t.Errorf("%s should be a valid locale: %s", test.localeID, err)
		}
		if !test.valid && err == nil {
			t.Errorf("%s should not be a valid locale", test.localeID)
		}
	}
}
