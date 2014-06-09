package language

import (
	"fmt"
	"strconv"
	"strings"
)

// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type operands struct {
	N float64 // absolute value of the source number (integer and decimals)
	I int64   // integer digits of n
	V int     // number of visible fraction digits in n, with trailing zeros
	W int     // number of visible fraction digits in n, without trailing zeros
	F int     // visible fractional digits in n, with trailing zeros
	T int     // visible fractional digits in n, without trailing zeros
}

func newOperands(v interface{}) (*operands, error) {
	switch v := v.(type) {
	case int:
		return newOperandsInt64(int64(v)), nil
	case int8:
		return newOperandsInt64(int64(v)), nil
	case int16:
		return newOperandsInt64(int64(v)), nil
	case int32:
		return newOperandsInt64(int64(v)), nil
	case int64:
		return newOperandsInt64(v), nil
	case string:
		return newOperandsString(v)
	case float32, float64:
		return nil, fmt.Errorf("floats should be formatted into a string")
	default:
		return nil, fmt.Errorf("invalid type %T; expected integer or string", v)
	}
}

func newOperandsInt64(i int64) *operands {
	if i < 0 {
		i = -i
	}
	return &operands{float64(i), i, 0, 0, 0, 0}
}

func newOperandsString(s string) (*operands, error) {
	if s[0] == '-' {
		s = s[1:]
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	ops := &operands{N: n}
	parts := strings.SplitN(s, ".", 2)
	ops.I, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}
	if len(parts) == 1 {
		return ops, nil
	}
	fraction := parts[1]
	ops.V = len(fraction)
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
		ops.F = int(f)
	}
	if ops.W > 0 {
		t, err := strconv.ParseInt(fraction[:ops.W], 10, 0)
		if err != nil {
			return nil, err
		}
		ops.T = int(t)
	}
	return ops, nil
}
