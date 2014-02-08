package plural

import (
	"testing"
)

func TestNewCategory(t *testing.T) {
	tests := []struct {
		src string
		cat Category
		err bool
	}{
		{"zero", Zero, false},
		{"one", One, false},
		{"two", Two, false},
		{"few", Few, false},
		{"many", Many, false},
		{"other", Other, false},
		{"asdf", Invalid, true},
	}

	for _, test := range tests {
		cat, err := NewCategory(test.src)
		wrongErr := (err != nil && !test.err) || (err == nil && test.err)
		if cat != test.cat || wrongErr {
			t.Errorf("New(%#v) returned %#v,%#v; expected %#v", test.src, cat, err, test.cat)
		}
	}

}
