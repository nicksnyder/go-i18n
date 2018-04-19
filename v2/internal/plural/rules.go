package plural

import "golang.org/x/text/language"

// Rules is a set of plural rules by language tag.
type Rules map[language.Tag]*Rule

// Rule returns the closest matching plural rule for the language tag
// or nil if no rule could be found.
func (r Rules) Rule(tag language.Tag) *Rule {
	for {
		if rule := r[tag]; rule != nil {
			return rule
		}
		tag = tag.Parent()
		if tag.IsRoot() {
			return nil
		}
	}
}
