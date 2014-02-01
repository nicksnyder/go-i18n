package i18n

import (
	"fmt"
)

// Language is a human language as defined by RFC 5646.
//
// Languages are identified by a 2 character language code
// optionally followed by a dash and a 4 character script subtag (e.g. en, zh-Hant)
//
// Languages implement CLDR plural rules as defined here:
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
type Language struct {
	ID               string
	Name             string
	PluralCategories map[PluralCategory]struct{}
	IntFunc          func(int64) PluralCategory
	FloatFunc        func(float64) PluralCategory
}

// Alphabetical by English name.
var languages = map[string]*Language{
	// Arabic
	"ar": &Language{
		ID:               "ar",
		Name:             "العربية",
		PluralCategories: newSet(Zero, One, Two, Few, Many, Other),
		IntFunc: func(i int64) PluralCategory {
			switch i {
			case 0:
				return Zero
			case 1:
				return One
			case 2:
				return Two
			default:
				mod100 := i % 100
				if mod100 >= 3 && mod100 <= 10 {
					return Few
				}
				if mod100 >= 11 {
					return Many
				}
				return Other
			}
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	// Chinese (Simplified)
	"zh-Hans": &Language{
		ID:               "zh-Hans",
		Name:             "汉语",
		PluralCategories: newSet(Other),
		IntFunc: func(i int64) PluralCategory {
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	// Chinese (Traditional)
	"zh-Hant": &Language{
		ID:               "zh-Hant",
		Name:             "漢語",
		PluralCategories: newSet(Other),
		IntFunc: func(i int64) PluralCategory {
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	"en": &Language{
		ID:               "en",
		Name:             "English",
		PluralCategories: newSet(One, Other),
		IntFunc: func(i int64) PluralCategory {
			if i == 1 {
				return One
			}
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	// French
	"fr": &Language{
		ID:               "fr",
		Name:             "Français",
		PluralCategories: newSet(One, Other),
		IntFunc: func(i int64) PluralCategory {
			if i == 0 || i == 1 {
				return One
			}
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			if f >= 0 && f < 2 {
				return One
			}
			return Other
		},
	},

	// German
	"de": &Language{
		ID:               "de",
		Name:             "Deutsch",
		PluralCategories: newSet(One, Other),
		IntFunc: func(i int64) PluralCategory {
			if i == 1 {
				return One
			}
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	// Italian
	"it": &Language{
		ID:               "it",
		Name:             "Italiano",
		PluralCategories: newSet(One, Other),
		IntFunc: func(i int64) PluralCategory {
			if i == 1 {
				return One
			}
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	// Japanese
	"ja": &Language{
		ID:               "ja",
		Name:             "日本語",
		PluralCategories: newSet(Other),
		IntFunc: func(i int64) PluralCategory {
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	// Spanish
	"es": &Language{
		ID:               "es",
		Name:             "Español",
		PluralCategories: newSet(One, Other),
		IntFunc: func(i int64) PluralCategory {
			if i == 1 {
				return One
			}
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},
}

// LanguageWithID returns the language identified by id
// or nil if the language is not registered.
func LanguageWithID(id string) *Language {
	return languages[id]
}

// RegisterLanguage adds l to the collection of available languages.
func RegisterLanguage(l *Language) {
	languages[l.ID] = l
}

func (l *Language) pluralCategory(count interface{}) (PluralCategory, error) {
	switch v := count.(type) {
	case int:
		return l.int64PluralCategory(int64(v)), nil
	case int8:
		return l.int64PluralCategory(int64(v)), nil
	case int16:
		return l.int64PluralCategory(int64(v)), nil
	case int32:
		return l.int64PluralCategory(int64(v)), nil
	case int64:
		return l.int64PluralCategory(v), nil
	case float32:
		return l.float64PluralCategory(float64(v)), nil
	case float64:
		return l.float64PluralCategory(v), nil
	default:
		return Invalid, fmt.Errorf("can't convert %#v to PluralCategory", v)
	}
}

func (l *Language) int64PluralCategory(i int64) PluralCategory {
	if i < 0 {
		i = -i
	}
	return l.IntFunc(i)
}

func (l *Language) float64PluralCategory(f float64) PluralCategory {
	if f < 0 {
		f = -f
	}
	if isInt64(f) {
		return l.IntFunc(int64(f))
	}
	return l.FloatFunc(f)
}

func isInt64(f float64) bool {
	return f == float64(int64(f))
}

func newSet(pluralCategories ...PluralCategory) map[PluralCategory]struct{} {
	set := make(map[PluralCategory]struct{}, len(pluralCategories))
	for _, pc := range pluralCategories {
		set[pc] = struct{}{}
	}
	return set
}
