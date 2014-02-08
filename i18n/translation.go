package i18n

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/plural"
)

// Translation is the interface that represents a translated string.
type Translation interface {
	MarshalInterface() interface{}
	ID() string
	Template(plural.Category) *template
	UntranslatedCopy() Translation
	Normalize(language *language.Language) Translation
	Backfill(src Translation) Translation
	Merge(Translation) Translation
	Incomplete(l *language.Language) bool
}

type byID []Translation

func (a byID) Len() int           { return len(a) }
func (a byID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byID) Less(i, j int) bool { return a[i].ID() < a[j].ID() }

// NewTranslation reflects on data to create a new Translation.
//
// data["id"] must be a string and data["translation"] must be either a string
// for a non-plural translation or a map[string]interface{} for a plural translation.
func NewTranslation(data map[string]interface{}) (Translation, error) {
	id, ok := data["id"].(string)
	if !ok {
		return nil, fmt.Errorf(`missing "id" key`)
	}
	switch translation := data["translation"].(type) {
	case string:
		tmpl, err := newTemplate(translation)
		if err != nil {
			return nil, err
		}
		return &singleTranslation{id, tmpl}, nil
	case map[string]interface{}:
		templates := make(map[plural.Category]*template, len(translation))
		for k, v := range translation {
			pc, err := plural.NewCategory(k)
			if err != nil {
				return nil, err
			}
			str, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf(`plural category "%s" has value of type %T; expected string`, pc, v)
			}
			tmpl, err := newTemplate(str)
			if err != nil {
				return nil, err
			}
			templates[pc] = tmpl
		}
		return &pluralTranslation{id, templates}, nil
	case nil:
		return nil, fmt.Errorf(`missing "translation" key`)
	default:
		return nil, fmt.Errorf(`unsupported type for "translation" key %T`, translation)
	}
}
