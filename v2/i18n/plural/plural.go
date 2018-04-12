package plural

import (
	"fmt"
	"strings"
)

// Form represents a language pluralization form as defined here:
// http://cldr.unicode.org/index/cldr-spec/plural-rules
type Form string

// All defined plural categories.
const (
	Invalid Form = ""
	Zero    Form = "zero"
	One     Form = "one"
	Two     Form = "two"
	Few     Form = "few"
	Many    Form = "many"
	Other   Form = "other"
)

// NewPluralForm returns src as a PluralForm
// or Invalid and a non-nil error if src is not a valid PluralForm.
func NewPluralForm(src string) (Form, error) {
	src = strings.ToLower(src)
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
	return Invalid, fmt.Errorf("invalid plural form %s", src)
}
