package plural

import (
	"testing"
)

func TestNewPluralForm(t *testing.T) {
	tests := []struct {
		src        string
		pluralForm PluralForm
		err        bool
	}{
		{"zero", Zero, false},
		{"Zero", Zero, false},
		{"one", One, false},
		{"One", One, false},
		{"two", Two, false},
		{"Two", Two, false},
		{"few", Few, false},
		{"Few", Few, false},
		{"many", Many, false},
		{"Many", Many, false},
		{"other", Other, false},
		{"Other", Other, false},
		{"", Invalid, true},
		{"asdf", Invalid, true},
	}
	for _, test := range tests {
		pluralForm, err := NewPluralForm(test.src)
		wrongErr := (err != nil && !test.err) || (err == nil && test.err)
		if pluralForm != test.pluralForm || wrongErr {
			t.Errorf("NewPlural(%#v) returned %#v,%#v; expected %#v", test.src, pluralForm, err, test.pluralForm)
		}
	}
}
