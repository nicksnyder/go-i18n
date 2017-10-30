package language

import (
	"strings"
	"sync"
)

// PluralSpec defines the CLDR plural rules for a language.
// http://www.unicode.org/cldr/charts/latest/supplemental/language_plural_rules.html
// http://unicode.org/reports/tr35/tr35-numbers.html#Operands
type PluralSpec struct {
	Plurals    map[Plural]struct{}
	PluralFunc func(*Operands) Plural
}

type pluralSpecs struct {
	specs map[string]*PluralSpec
	sync.RWMutex
}

func (ps *pluralSpecs) get(id string) *PluralSpec {
	ps.RLock()
	defer ps.RUnlock()
	return ps.specs[id]
}

func (ps *pluralSpecs) set(id string, p *PluralSpec) {
	ps.Lock()
	ps.specs[id] = p
	ps.Unlock()
}

var pluralSpecsStore = &pluralSpecs{specs: make(map[string]*PluralSpec)}

func normalizePluralSpecID(id string) string {
	id = strings.Replace(id, "_", "-", -1)
	id = strings.ToLower(id)
	return id
}

// RegisterPluralSpec registers a new plural spec for the language ids.
func RegisterPluralSpec(ids []string, ps *PluralSpec) {
	for _, id := range ids {
		id = normalizePluralSpecID(id)
		pluralSpecsStore.set(id, ps)
	}
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
		if spec := pluralSpecsStore.get(subtag); spec != nil {
			return spec
		}
		end := strings.LastIndex(subtag, "-")
		if end == -1 {
			return nil
		}
		subtag = subtag[:end]
	}
}

func newPluralSet(plurals ...Plural) map[Plural]struct{} {
	set := make(map[Plural]struct{}, len(plurals))
	for _, plural := range plurals {
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
