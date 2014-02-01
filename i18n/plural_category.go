package i18n

import (
	"fmt"
)

// PluralCategory represents a language pluralization form as defined here:
// http://cldr.unicode.org/index/cldr-spec/plural-rules
type PluralCategory string

// Enumeration of all valid PluralCategory
const (
	Invalid PluralCategory = "invalid"
	Zero                   = "zero"
	One                    = "one"
	Two                    = "two"
	Few                    = "few"
	Many                   = "many"
	Other                  = "other"
)

// NewPluralCategory returns converts src to a PluralCategory
// or returns Invalid and a non-nil error if src is not a valid PluralCategory.
func NewPluralCategory(src string) (PluralCategory, error) {
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
