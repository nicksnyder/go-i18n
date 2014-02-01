package i18n

import (
	"fmt"
)

type Language struct {
	Code             string
	Name             string
	PluralCategories map[PluralCategory]struct{}
	IntFunc          func(int64) PluralCategory
	FloatFunc        func(float64) PluralCategory
}

// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// Alphabetical by English name.
var languages = map[string]*Language{
	"ar": &Language{
		Code:             "ar",
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
		Code:             "zh-Hans",
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
		Code:             "zh-Hant",
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
		Code:             "en",
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

	"fr": &Language{
		Code:             "fr",
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

	"de": &Language{
		Code:             "de",
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

	"it": &Language{
		Code:             "it",
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

	"ja": &Language{
		Code:             "ja",
		Name:             "日本語",
		PluralCategories: newSet(Other),
		IntFunc: func(i int64) PluralCategory {
			return Other
		},
		FloatFunc: func(f float64) PluralCategory {
			return Other
		},
	},

	"es": &Language{
		Code:             "es",
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

func LanguageWithCode(code string) *Language {
	return languages[code]
}

func RegisterLanguage(l *Language) {
	languages[l.Code] = l
}

func (l *Language) String() string {
	return l.Name
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
