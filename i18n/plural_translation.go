package i18n

type pluralTranslation struct {
	id        string
	templates map[PluralCategory]*template
}

func (pt *pluralTranslation) MarshalInterface() interface{} {
	return map[string]interface{}{
		"id":          pt.id,
		"translation": pt.templates,
	}
}

func (pt *pluralTranslation) Id() string {
	return pt.id
}

func (pt *pluralTranslation) Template(pc PluralCategory) *template {
	return pt.templates[pc]
}

func (pt *pluralTranslation) UntranslatedCopy() Translation {
	return &pluralTranslation{pt.id, make(map[PluralCategory]*template)}
}

func (pt *pluralTranslation) Normalize(l *Language) Translation {
	// Delete plural categories that don't belong to this language.
	for pc := range pt.templates {
		if _, ok := l.PluralCategories[pc]; !ok {
			delete(pt.templates, pc)
		}
	}
	// Create map entries for missing valid categories.
	for pc := range l.PluralCategories {
		if _, ok := pt.templates[pc]; !ok {
			pt.templates[pc] = mustNewTemplate("")
		}
	}
	return pt
}

func (pt *pluralTranslation) Backfill(src Translation) Translation {
	for pc, t := range pt.templates {
		if t == nil || t.src == "" {
			pt.templates[pc] = src.Template(Other)
		}
	}
	return pt
}

func (pt *pluralTranslation) Merge(t Translation) Translation {
	other, ok := t.(*pluralTranslation)
	if !ok || pt.Id() != t.Id() {
		return t
	}
	for pluralCategory, template := range other.templates {
		if template != nil && template.src != "" {
			pt.templates[pluralCategory] = template
		}
	}
	return pt
}

func (pt *pluralTranslation) Incomplete(l *Language) bool {
	for pc := range l.PluralCategories {
		if t := pt.templates[pc]; t == nil || t.src == "" {
			return true
		}
	}
	return false
}

var _ = Translation(&pluralTranslation{})
