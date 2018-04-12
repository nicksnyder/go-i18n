package plural

import "golang.org/x/text/language"

// PluralRule defines the CLDR plural rules for a language.
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type PluralRule struct {
	PluralForms    map[PluralForm]struct{}
	PluralFormFunc func(*Operands) PluralForm
}

func addPluralRules(rules map[language.Base]*PluralRule, ids []string, ps *PluralRule) {
	for _, id := range ids {
		if id == "root" {
			continue
		}
		base := language.MustParseBase(id)
		rules[base] = ps
	}
}

func newPluralFormSet(pluralForms ...PluralForm) map[PluralForm]struct{} {
	set := make(map[PluralForm]struct{}, len(pluralForms))
	for _, plural := range pluralForms {
		set[plural] = struct{}{}
	}
	return set
}

func intInRange(i, from, to int64) bool {
	return from <= i && i <= to
}

func intEqualsAny(i int64, any ...int64) bool {
	for _, a := range any {
		if i == a {
			return true
		}
	}
	return false
}
