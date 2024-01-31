package plural

import (
	"fmt"
	"strconv"
	"strings"
)

// Operands is a representation of http://unicode.org/reports/tr35/tr35-numbers.html#Operands
// If there is a compact decimal exponent value C, then the N, I, V, W, F, and T values are computed after shifting the decimal point in the original by the ‘c’ value.
// So for 1.2c3, the values are the same as those of 1200: i=1200 and f=0.
// Similarly, 1.2005c3 has i=1200 and f=5 (corresponding to 1200.5).
type Operands struct {
	N float64 // absolute value of the source number (integer and decimals)
	I int64   // integer digits of n
	V int64   // number of visible fraction digits in n, with trailing zeros
	W int64   // number of visible fraction digits in n, without trailing zeros
	F int64   // visible fractional digits in n, with trailing zeros
	T int64   // visible fractional digits in n, without trailing zeros
	C int64   // compact decimal exponent value: exponent of the power of 10 used in compact decimal formatting.
}

// NEqualsAny returns true if o represents an integer equal to any of the arguments.
func (o *Operands) NEqualsAny(any ...int64) bool {
	for _, i := range any {
		if o.I == i && o.T == 0 {
			return true
		}
	}
	return false
}

// NModEqualsAny returns true if o represents an integer equal to any of the arguments modulo mod.
func (o *Operands) NModEqualsAny(mod int64, any ...int64) bool {
	modI := o.I % mod
	for _, i := range any {
		if modI == i && o.T == 0 {
			return true
		}
	}
	return false
}

// NInRange returns true if o represents an integer in the closed interval [from, to].
func (o *Operands) NInRange(from, to int64) bool {
	return o.T == 0 && from <= o.I && o.I <= to
}

// NModInRange returns true if o represents an integer in the closed interval [from, to] modulo mod.
func (o *Operands) NModInRange(mod, from, to int64) bool {
	modI := o.I % mod
	return o.T == 0 && from <= modI && modI <= to
}

// NewOperands returns the operands for number.
func NewOperands(number interface{}) (*Operands, error) {
	switch number := number.(type) {
	case int:
		return newOperandsInt64(int64(number)), nil
	case int8:
		return newOperandsInt64(int64(number)), nil
	case int16:
		return newOperandsInt64(int64(number)), nil
	case int32:
		return newOperandsInt64(int64(number)), nil
	case int64:
		return newOperandsInt64(number), nil
	case string:
		return newOperandsString(number)
	case float32, float64:
		return nil, fmt.Errorf("floats should be formatted into a string")
	default:
		return nil, fmt.Errorf("invalid type %T; expected integer or string", number)
	}
}

func newOperandsInt64(i int64) *Operands {
	if i < 0 {
		i = -i
	}
	return &Operands{float64(i), i, 0, 0, 0, 0, 0}
}

func splitSignificandExponent(s string) (significand, exponent string) {
	i := strings.IndexAny(s, "eE")
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func shiftDecimalLeft(s string, n int) string {
	if n <= 0 {
		return s
	}
	i := strings.IndexRune(s, '.')
	tilt := 0
	if i < 0 {
		i = len(s)
		tilt = -1
	}
	switch {
	case n == i:
		return "0." + s[:i] + s[i+1+tilt:]
	case n > i:
		return "0." + strings.Repeat("0", n-i) + s[:i] + s[i+1+tilt:]
	default:
		return s[:i-n] + "." + s[i-n:i] + s[i+1+tilt:]
	}
}

func shiftDecimalRight(s string, n int) string {
	if n <= 0 {
		return s
	}
	i := strings.IndexRune(s, '.')
	if i < 0 {
		return s + strings.Repeat("0", n)
	}
	switch rest := len(s) - i - 1; {
	case n == rest:
		return s[:i] + s[i+1:]
	case n > rest:
		return s[:i] + s[i+1:] + strings.Repeat("0", n-rest)
	default:
		return s[:i] + s[i+1:i+1+n] + "." + s[i+1+n:]
	}
}

func applyExponent(s string, exponent int) string {
	switch {
	case exponent > 0:
		return shiftDecimalRight(s, exponent)
	case exponent < 0:
		return shiftDecimalLeft(s, -exponent)
	}
	return s
}

func newOperandsString(s string) (*Operands, error) {
	if s[0] == '-' {
		s = s[1:]
	}
	ops := &Operands{}
	var err error
	ops.N, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	significand, exponent := splitSignificandExponent(s)
	if exponent != "" {
		// We are storing C as an int64 but only allowing
		// numbers that fit into the bitsize of an int
		// so C is safe to cast as a int later.
		ops.C, err = strconv.ParseInt(exponent, 10, 0)
		if err != nil {
			return nil, err
		}
	}
	value := applyExponent(significand, int(ops.C))
	parts := strings.SplitN(value, ".", 2)
	ops.I, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	if len(parts) == 1 {
		return ops, nil
	}
	fraction := parts[1]
	ops.V = int64(len(fraction))
	for i := ops.V - 1; i >= 0; i-- {
		if fraction[i] != '0' {
			ops.W = i + 1
			break
		}
	}
	if ops.V > 0 {
		f, err := strconv.ParseInt(fraction, 10, 0)
		if err != nil {
			return nil, err
		}
		ops.F = f
	}
	if ops.W > 0 {
		t, err := strconv.ParseInt(fraction[:ops.W], 10, 0)
		if err != nil {
			return nil, err
		}
		ops.T = t
	}
	return ops, nil
}
