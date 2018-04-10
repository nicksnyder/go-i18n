package i18n

import (
	"fmt"
	"strings"
)

// PluralForm represents a language pluralization form as defined here:
// http://cldr.unicode.org/index/cldr-spec/plural-rules
type PluralForm string

// All defined plural categories.
const (
	Invalid PluralForm = ""
	Zero    PluralForm = "zero"
	One     PluralForm = "one"
	Two     PluralForm = "two"
	Few     PluralForm = "few"
	Many    PluralForm = "many"
	Other   PluralForm = "other"
)

// NewPluralForm returns src as a PluralForm
// or Invalid and a non-nil error if src is not a valid PluralForm.
func NewPluralForm(src string) (PluralForm, error) {
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
