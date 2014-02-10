// Package language defines languages that implement CLDR pluralization.
package language

import (
	"github.com/nicksnyder/go-i18n/i18n/plural"
)

// Language is a written human language.
//
// A Language is identified by a 2 character language code
// optionally followed by a dash and a 4 character script subtag (e.g. en, zh-Hant)
// as defined by RFC 5646.
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
		Name:             "العربية",
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
		Name:             "Català",
		PluralCategories: newSet(plural.One, plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			if ops.I == 1 && ops.V == 0 {
				return plural.One
			}
			return plural.Other
		},
	},

	// Chinese (Simplified)
	// TODO: inconsistent with pt-BR
	"zh-Hans": &Language{
		ID:               "zh-Hans",
		Name:             "汉语",
		PluralCategories: newSet(plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			return plural.Other
		},
	},

	// Chinese (Traditional)
	// TODO: inconsistent with pt-BR
	"zh-Hant": &Language{
		ID:               "zh-Hant",
		Name:             "漢語",
		PluralCategories: newSet(plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			return plural.Other
		},
	},

	// Czech
	"cs": &Language{
		ID:               "cs",
		Name:             "Čeština",
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
		Name:             "Dansk",
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
		Name:             "Nederlands",
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
		Name:             "English",
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
		Name:             "Français",
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
		Name:             "Deutsch",
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
		Name:             "Italiano",
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
		Name:             "日本語",
		PluralCategories: newSet(plural.Other),
		PluralFunc: func(ops *plural.Operands) plural.Category {
			return plural.Other
		},
	},

	// Portuguese (European)
	"pt": &Language{
		ID:               "pt",
		Name:             "Português",
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
		Name:             "Português Brasileiro",
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
		Name:             "Español",
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

func newSet(pluralCategories ...plural.Category) map[plural.Category]struct{} {
	set := make(map[plural.Category]struct{}, len(pluralCategories))
	for _, pc := range pluralCategories {
		set[pc] = struct{}{}
	}
	return set
}
