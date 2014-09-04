package language

import (
	"strings"
	"math"
)

// PluralSpec defines the CLDR plural rules for a language.
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type PluralSpec struct {
	Plurals    map[Plural]struct{}
	PluralFunc func(*operands) Plural
}

// Plural returns the plural category for number as defined by
// the language's CLDR plural rules.
func (ps *PluralSpec) Plural(number interface{}) (Plural, error) {
	ops, err := newOperands(number)
	if err != nil {
		return Invalid, err
	}
	return ps.PluralFunc(ops), nil
}

// getPluralSpec returns the PluralSpec that matches the longest prefix of tag.
// It returns nil if no PluralSpec matches tag.
func getPluralSpec(tag string) *PluralSpec {
	tag = NormalizeTag(tag)
	subtag := tag
	for {
		if spec := pluralSpecs[subtag]; spec != nil {
			return spec
		}
		end := strings.LastIndex(subtag, "-")
		if end == -1 {
			return nil
		}
		subtag = subtag[:end]
	}
}

// Alphabetical by English name.
var pluralSpecs = map[string]*PluralSpec{
	// Arabic
	"ar": &PluralSpec{
		Plurals: newPluralSet(Zero, One, Two, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.W == 0 {
				switch ops.I {
				case 0:
					return Zero
				case 1:
					return One
				case 2:
					return Two
				default:
					mod100 := ops.I % 100
					if mod100 >= 3 && mod100 <= 10 {
						return Few
					}
					if mod100 >= 11 {
						return Many
					}
				}
			}
			return Other
		},
	},

	// Belarusian
	"be": &PluralSpec{
		Plurals: newPluralSet(One, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			mod10 := math.Mod(ops.N, 10)
			mod100 := math.Mod(ops.N, 100)
			if ops.T == 0 && mod10 == 1 && mod100 != 11 {
				return One
			}
			if ops.T == 0 && mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
				return Few
			}
			if	(ops.T == 0 && mod10 == 0) ||
				(ops.T == 0 && mod10 >= 5 && mod10 <= 9) ||
				(ops.T == 0 && mod100 >= 11 && mod100 <= 14) {
				return Many
			}
			return Other
		},
	},

	// Catalan
	"ca": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// Chinese
	// There is no need to distinguish between simplified and traditional
	// since they have the same pluralization.
	"zh": &PluralSpec{
		Plurals: newPluralSet(Other),
		PluralFunc: func(ops *operands) Plural {
			return Other
		},
	},

	// Czech
	"cs": &PluralSpec{
		Plurals: newPluralSet(One, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			if ops.I >= 2 && ops.I <= 4 && ops.V == 0 {
				return Few
			}
			if ops.V > 0 {
				return Many
			}
			return Other
		},
	},

	// Danish
	"da": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 || (ops.I == 0 && ops.T != 0) {
				return One
			}
			return Other
		},
	},

	// Dutch
	"nl": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// English
	"en": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// French
	"fr": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 0 || ops.I == 1 {
				return One
			}
			return Other
		},
	},

	// German
	"de": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// Icelandic
	"is": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if (ops.T == 0 && ops.I % 10 == 1 && ops.I % 100 != 11) || ops.T != 0 {
				return One
			}
			return Other
		},
	},

	// Italian
	"it": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// Japanese
	"ja": &PluralSpec{
		Plurals: newPluralSet(Other),
		PluralFunc: func(ops *operands) Plural {
			return Other
		},
	},

	// Lithuanian
	"lt": &PluralSpec{
		Plurals: newPluralSet(One, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.F != 0 {
				return Many
			}
			mod100 := ops.I % 100
			if mod100 < 11 || mod100 > 19 {
				switch ops.I % 10 {
				case 0:
					return Other
				case 1:
					return One
				default:
					return Few
				}
			}
			return Other
		},
	},

	// Polish
	"pl": &PluralSpec{
		Plurals: newPluralSet(One, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.V == 0 && ops.I == 1 {
				return One
			}
			mod10 := ops.I % 10
			mod100 := ops.I % 100
			if ops.V == 0 && mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
				return Few
			}
			if	(ops.V == 0 && ops.I != 1 && mod10 >= 0 && mod10 <= 1) ||
				(ops.V == 0 && mod10 >= 5 && mod10 <= 9) ||
				(ops.V == 0 && mod100 >= 12 && mod100 <= 14) {
				return Many
			}
			return Other
		},
	},
	
	// Portuguese (European)
	"pt": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// Portuguese (Brazilian)
	"pt-br": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if (ops.I == 1 && ops.V == 0) || (ops.I == 0 && ops.T == 1) {
				return One
			}
			return Other
		},
	},

	// Russian
	"ru": &PluralSpec{
		Plurals: newPluralSet(One, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			mod10 := ops.I % 10
			mod100 := ops.I % 100
			if ops.V == 0 && mod10 == 1 && mod100 != 11 {
				return One
			}
			if ops.V == 0 && mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
				return Few
			}
			if	(ops.V == 0 && mod10 == 0) ||
				(ops.V == 0 && mod10 >= 5 && mod10 <= 9) ||
				(ops.V == 0 && mod100 >= 11 && mod100 <= 14) {
				return Many
			}
			return Other
		},
	},

	// Spanish
	"es": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.W == 0 {
				return One
			}
			return Other
		},
	},

	// Bulgarian
	"bg": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.W == 0 {
				return One
			}
			return Other
		},
	},

	// Swedish
	"sv": &PluralSpec{
		Plurals: newPluralSet(One, Other),
		PluralFunc: func(ops *operands) Plural {
			if ops.I == 1 && ops.V == 0 {
				return One
			}
			return Other
		},
	},

	// Ukrainian
	"uk": &PluralSpec{
		Plurals: newPluralSet(One, Few, Many, Other),
		PluralFunc: func(ops *operands) Plural {
			mod10 := ops.I % 10
			mod100 := ops.I % 100
			if ops.V == 0 && mod10 == 1 && mod100 != 11 {
				return One
			}
			if ops.V == 0 && mod10 >= 2 && mod10 <= 4 && !(mod100 >= 12 && mod100 <= 14) {
				return Few
			}
			if	(ops.V == 0 && mod10 == 0) ||
				(ops.V == 0 && mod10 >= 5 && mod10 <= 9) ||
				(ops.V == 0 && mod100 >= 11 && mod100 <= 14) {
				return Many
			}
			return Other
		},
	},
}

func newPluralSet(plurals ...Plural) map[Plural]struct{} {
	set := make(map[Plural]struct{}, len(plurals))
	for _, plural := range plurals {
		set[plural] = struct{}{}
	}
	return set
}
