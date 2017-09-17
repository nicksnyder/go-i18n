package language

import (
	"reflect"
	"testing"
)

func TestNewOperands(t *testing.T) {
	tests := []struct {
		input interface{}
		op    *Operand
		err   bool
	}{
		{int64(0), &Operand{0.0, 0, 0, 0, 0, 0}, false},
		{int64(1), &Operand{1.0, 1, 0, 0, 0, 0}, false},
		{"0", &Operand{0.0, 0, 0, 0, 0, 0}, false},
		{"1", &Operand{1.0, 1, 0, 0, 0, 0}, false},
		{"1.0", &Operand{1.0, 1, 1, 0, 0, 0}, false},
		{"1.00", &Operand{1.0, 1, 2, 0, 0, 0}, false},
		{"1.3", &Operand{1.3, 1, 1, 1, 3, 3}, false},
		{"1.30", &Operand{1.3, 1, 2, 1, 30, 3}, false},
		{"1.03", &Operand{1.03, 1, 2, 2, 3, 3}, false},
		{"1.230", &Operand{1.23, 1, 3, 2, 230, 23}, false},
		{"20.0230", &Operand{20.023, 20, 4, 3, 230, 23}, false},
		{20.0230, nil, true},
	}
	for _, test := range tests {
		op, err := newOperand(test.input)
		if err != nil && !test.err {
			t.Errorf("newOperands(%#v) unexpected error: %s", test.input, err)
		} else if err == nil && test.err {
			t.Errorf("newOperands(%#v) returned %#v; expected error", test.input, op)
		} else if !reflect.DeepEqual(op, test.op) {
			t.Errorf("newOperands(%#v) returned %#v; expected %#v", test.input, op, test.op)
		}
	}
}

func BenchmarkNewOperand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := newOperand("1234.56780000"); err != nil {
			b.Fatal(err)
		}
	}
}
