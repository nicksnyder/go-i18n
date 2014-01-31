package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"launchpad.net/goyaml"
	"path/filepath"
	"reflect"
	"sort"
)

type marshalFunc func(interface{}) ([]byte, error)

func Merge(translationFiles []string, sourceLocaleID, outdir, format string) error {
	if len(translationFiles) < 1 {
		return fmt.Errorf("need at least one translation file to parse")
	}

	if _, err := NewLocale(sourceLocaleID); err != nil {
		return fmt.Errorf("invalid source locale %s: %s", sourceLocaleID, err)
	}

	marshal, err := newMarshalFunc(format)
	if err != nil {
		return err
	}

	bundle := NewBundle()
	for _, tf := range translationFiles {
		if err := bundle.LoadTranslationFile(tf); err != nil {
			return fmt.Errorf("failed to load translation file %s because %s\n", tf, err)
		}
	}

	sourceTranslations := bundle.translations[sourceLocaleID]
	for translationID, src := range sourceTranslations {
		for _, localeTranslations := range bundle.translations {
			if dst := localeTranslations[translationID]; dst == nil || reflect.TypeOf(src) != reflect.TypeOf(dst) {
				localeTranslations[translationID] = src.UntranslatedCopy()
			}
		}
	}

	for localeID, localeTranslations := range bundle.translations {
		writeFile := writeFileFunc(outdir, localeID, format, marshal)
		locale := mustNewLocale(localeID)
		all := filter(localeTranslations, func(t Translation) Translation {
			return t.Normalize(locale.Language)
		})
		if err := writeFile("all", all); err != nil {
			return err
		}

		untranslated := filter(localeTranslations, func(t Translation) Translation {
			if t.Incomplete(locale.Language) {
				return t.Normalize(locale.Language).Backfill(sourceTranslations[t.ID()])
			}
			return nil
		})
		if err := writeFile("untranslated", untranslated); err != nil {
			return err
		}
	}
	return nil
}

func filter(translations map[string]Translation, filter func(Translation) Translation) []Translation {
	filtered := make([]Translation, 0, len(translations))
	for _, translation := range translations {
		if t := filter(translation); t != nil {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func newMarshalFunc(format string) (marshalFunc, error) {
	switch format {
	case "json":
		return func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		}, nil
		/*
			case "yaml":
				return func(v interface{}) ([]byte, error) {
					return goyaml.Marshal(v)
				}, nil
		*/
	}
	return nil, fmt.Errorf("unsupported format: %s\n", format)
}

func marshalInterface(translations []Translation) []interface{} {
	mi := make([]interface{}, len(translations))
	for i, translation := range translations {
		mi[i] = translation.MarshalInterface()
	}
	return mi
}

func writeFileFunc(outdir, localeID, format string, marshal marshalFunc) func(string, []Translation) error {
	return func(label string, translations []Translation) error {
		sort.Sort(byID(translations))
		buf, err := marshal(marshalInterface(translations))
		if err != nil {
			return fmt.Errorf("failed to marshal %s strings to %s because %s", localeID, format, err)
		}
		filename := filepath.Join(outdir, fmt.Sprintf("%s.%s.%s", localeID, label, format))
		if err := ioutil.WriteFile(filename, buf, 0666); err != nil {
			return fmt.Errorf("failed to write %s because %s", filename, err)
		}
		return nil
	}
}
