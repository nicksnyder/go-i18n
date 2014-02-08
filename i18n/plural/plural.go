// Package plural defines CLDR plural categories.
package plural

import (
	"fmt"
)

// Category represents a language pluralization form as defined here:
// http://cldr.unicode.org/index/cldr-spec/plural-rules
type Category string

// All defined plural categories.
const (
	Invalid Category = "invalid"
	Zero             = "zero"
	One              = "one"
	Two              = "two"
	Few              = "few"
	Many             = "many"
	Other            = "other"
)

// NewCategory returns src as a Category
// or Invalid and a non-nil error if src is not a valid Category.
func NewCategory(src string) (Category, error) {
	switch src {
	case "zero":
		return Zero, nil
	case "one":
		return One, nil
	case "two":
		return Two, nil
	case "few":
		return Few, nil
	case "many":
		return Many, nil
	case "other":
		return Other, nil
	}
	return Invalid, fmt.Errorf("invalid plural category %s", src)
}
