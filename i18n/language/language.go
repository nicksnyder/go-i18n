// Package language defines languages that implement CLDR pluralization.
package language

import (
	"github.com/nicksnyder/go-i18n/i18n/plural"
)

// Language is a written human language.
//
// Languages are identified by tags defined by RFC 5646.
//
// Typically language tags are a 2 character language code (ISO 639-1)
// optionally followed by a dash and a 2 character country code (ISO 3166-1).
// (e.g. en, pt-BR)
//
// A Language implements CLDR plural rules as defined here:
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type Language struct {
	ID               string
	Name             string
	PluralCategories map[plural.Category]struct{}
	PluralFunc       func(*plural.Operands) plural.Category
}

// Alphabetical by English name.
var languages = map[string]*Language{
	// Arabic
	"ar": &Language{
		ID:               "ar",
		PluralCategories: newSet(plural.Zero, plural.One, plural.Two, plural.Few, plural.Many, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.W == 0 {
				// Integer case
				switch ops.I {
				case 0:
					return plural.Zero
				case 1:
					return plural.One
				case 2:
					return plural.Two
				default:
					mod100 := ops.I % 100
					if mod100 >= 3 && mod100 <= 10 {
						return plural.Few
					}
					if mod100 >= 11 {
						return plural.Many
					}
				}
			}
			return plural.Other
		},
	},

	// Catalan
	"ca": &Language{
		ID:               "ca",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// Chinese
	// There is no need to distinguish between simplified and traditional
	// since they have the same pluralization.
	"zh": &Language{
		ID:               "zh",
		PluralCategories: newSet(plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			return plural.Other
		},
	},

	// Czech
	"cs": &Language{
		ID:               "cs",
		PluralCategories: newSet(plural.One, plural.Few, plural.Many, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			if ops.I >= 2 && ops.I <= 4 && ops.V == 0 {
				return plural.Few
			}
			if ops.V > 0 {
				return plural.Many
			}
			return plural.Other
		},
	},

	// Danish
	"da": &Language{
		ID:               "da",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 || (ops.I == 0 && ops.T != 0) {
				return plural.One
			}
			return plural.Other
		},
	},

	// Dutch
	"nl": &Language{
		ID:               "nl",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// English
	"en": &Language{
		ID:               "en",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// French
	"fr": &Language{
		ID:               "fr",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 0 || ops.I == 1 {
				return plural.One
			}
			return plural.Other
		},
	},

	// German
	"de": &Language{
		ID:               "de",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// Italian
	"it": &Language{
		ID:               "it",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// Japanese
	"ja": &Language{
		ID:               "ja",
		PluralCategories: newSet(plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			return plural.Other
		},
	},

	// Lithuanian
	"lt": &Language{
		ID:               "lt",
		PluralCategories: newSet(plural.One, plural.Few, plural.Many, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.F != 0 {
				return plural.Many
			}
			mod100 := ops.I % 100
			if mod100 < 11 || mod100 > 19 {
				switch ops.I % 10 {
				case 0:
					return plural.Other
				case 1:
					return plural.One
				default:
					return plural.Few
				}
			}
			return plural.Other
		},
	},

	// Portuguese (European)
	"pt": &Language{
		ID:               "pt",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// Portuguese (Brazilian)
	"pt-BR": &Language{
		ID:               "pt-BR",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if (ops.I == 1 && ops.V == 0) || (ops.I == 0 && ops.T == 1) {
				return plural.One
			}
			return plural.Other
		},
	},

	// Spanish
	"es": &Language{
		ID:               "es",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.W == 0 {
				return plural.One
			}
			return plural.Other
		},
	},
}

// LanguageWithID returns the language identified by id
// or nil if the language is not registered.
func LanguageWithID(id string) *Language {
	return languages[id]
}

// Register adds Language l to the collection of available languages.
func Register(l *Language) {
	languages[l.ID] = l
}

// PluralCategory returns the plural category for number as defined by
// the language's CLDR plural rules.
func (l *Language) PluralCategory(number interface{}) (plural.Category, error) {
	ops, err := plural.NewOperands(number)
	if err != nil {
		return plural.Invalid, err
	}
	return l.PluralFunc(ops), nil
}

func (l *Language) String() string {
	return l.ID
}

func newSet(pluralCategories ...plural.Category) map[plural.Category]struct{} {
	set := make(map[plural.Category]struct{}, len(pluralCategories))
	for _, pc := range pluralCategories {
		set[pc] = struct{}{}
	}
	return set
}
