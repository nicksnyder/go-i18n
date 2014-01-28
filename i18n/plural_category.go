package i18n

import (
	"fmt"
)

type PluralCategory string

const (
	Invalid PluralCategory = "invalid"
	Zero                   = "zero"
	One                    = "one"
	Two                    = "two"
	Few                    = "few"
	Many                   = "many"
	Other                  = "other"
)

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
