package plural

import (
	"reflect"
	"testing"
)

func TestNewOperands(t *testing.T) {
	tests := []struct {
		input interface{}
		ops   *Operands
		err   bool
	}{
		{int64(0), &Operands{0.0, 0, 0, 0, 0, 0, 0}, false},
		{int64(1), &Operands{1.0, 1, 0, 0, 0, 0, 0}, false},
		{"0", &Operands{0.0, 0, 0, 0, 0, 0, 0}, false},
		{"1", &Operands{1.0, 1, 0, 0, 0, 0, 0}, false},
		{"1.0", &Operands{1.0, 1, 1, 0, 0, 0, 0}, false},
		{"1.00", &Operands{1.0, 1, 2, 0, 0, 0, 0}, false},
		{"1.3", &Operands{1.3, 1, 1, 1, 3, 3, 0}, false},
		{"1.30", &Operands{1.3, 1, 2, 1, 30, 3, 0}, false},
		{"1.03", &Operands{1.03, 1, 2, 2, 3, 3, 0}, false},
		{"1.230", &Operands{1.23, 1, 3, 2, 230, 23, 0}, false},
		{"20.0230", &Operands{20.023, 20, 4, 3, 230, 23, 0}, false},
		{20.0230, nil, true},

		{"1200", &Operands{1200, 1200, 0, 0, 0, 0, 0}, false},
		{"1.2e3", &Operands{1200, 1200, 0, 0, 0, 0, 3}, false},
		{"1.2E3", &Operands{1200, 1200, 0, 0, 0, 0, 3}, false},

		{"1234", &Operands{1234, 1234, 0, 0, 0, 0, 0}, false},
		{"1234e0", &Operands{1234, 1234, 0, 0, 0, 0, 0}, false},
		{"123.4e1", &Operands{1234, 1234, 0, 0, 0, 0, 1}, false},
		{"12.34e2", &Operands{1234, 1234, 0, 0, 0, 0, 2}, false},
		{"1.234e3", &Operands{1234, 1234, 0, 0, 0, 0, 3}, false},
		{"0.1234e4", &Operands{1234, 1234, 0, 0, 0, 0, 4}, false},
		{"0.01234e5", &Operands{1234, 1234, 0, 0, 0, 0, 5}, false},

		{"1234.0", &Operands{1234, 1234, 1, 0, 0, 0, 0}, false},
		{"12340e-1", &Operands{1234, 1234, 1, 0, 0, 0, -1}, false},

		{"1200.5", &Operands{1200.5, 1200, 1, 1, 5, 5, 0}, false},
		{"1.2005e3", &Operands{1200.5, 1200, 1, 1, 5, 5, 3}, false},

		{"1200e3", &Operands{1200000, 1200000, 0, 0, 0, 0, 3}, false},

		{"0.0012340", &Operands{0.001234, 0, 7, 6, 12340, 1234, 0}, false},
		{"0.012340e-1", &Operands{0.001234, 0, 7, 6, 12340, 1234, -1}, false},
		{"0.12340e-2", &Operands{0.001234, 0, 7, 6, 12340, 1234, -2}, false},
		{"1.2340e-3", &Operands{0.001234, 0, 7, 6, 12340, 1234, -3}, false},
		{"12.340e-4", &Operands{0.001234, 0, 7, 6, 12340, 1234, -4}, false},
		{"123.40e-5", &Operands{0.001234, 0, 7, 6, 12340, 1234, -5}, false},
		{"1234.0e-6", &Operands{0.001234, 0, 7, 6, 12340, 1234, -6}, false},
		{"12340e-7", &Operands{0.001234, 0, 7, 6, 12340, 1234, -7}, false},
	}
	for _, test := range tests {
		ops, err := NewOperands(test.input)
		if err != nil && !test.err {
			t.Errorf("NewOperands(%#v) unexpected error: %s", test.input, err)
		} else if err == nil && test.err {
			t.Errorf("NewOperands(%#v) returned %#v; expected error", test.input, ops)
		} else if !reflect.DeepEqual(ops, test.ops) {
			t.Errorf("NewOperands(%#v) returned %#v; expected %#v", test.input, ops, test.ops)
		}
	}
}

func BenchmarkNewOperand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := NewOperands("1234.56780000"); err != nil {
			b.Fatal(err)
		}
	}
}
