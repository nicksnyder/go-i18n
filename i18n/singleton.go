package i18n

import (
	"sync"
)

var (
	localizerSingletonOnce sync.Once
	localizerInstance      *Localizer
	localizerMutex         = &sync.Mutex{}

	notLocalized string

	abcParams []string
)

func GetLocalizerInstance() *Localizer {
	localizerSingletonOnce.Do(func() {
		if localizerInstance == nil {
			localizerInstance = new(Localizer)
		}
	})
	return localizerInstance
}

// SetLocalizerInstance must be run before using `GetLocalizerInstance()` in a multi-threading manner.
func SetLocalizerInstance(l *Localizer) {
	if localizerMutex != nil {
		localizerMutex.Lock()
		defer localizerMutex.Unlock()
	}

	localizerInstance = l
}

// SetUseNotLocalizedInfo is optional.
func SetUseNotLocalizedInfo(use bool) {
	if use {
		notLocalized = Localize("NotLocalized")
	} else {
		notLocalized = ""
	}
}

func ResetSingletonContext() {
	SetUseNotLocalizedInfo(false)
	SetABCParams(nil)
}

// Localize loads localization for static text.
//
// Use only if you called SetLocalizerInstance().
func Localize(id string) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{MessageID: id})
	if err != nil {
		return notLocalized
	}
	return localization
}

// LocalizePlural loads localization for text with plural count.
//
// Use only if you called SetLocalizerInstance().
func LocalizePlural(id string, count int) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{
		MessageID:   id,
		PluralCount: count,
	})
	if err != nil {
		return notLocalized
	}
	return localization
}

// LocalizeTemplate loads localization using template.
//
// Use only if you called SetLocalizerInstance().
func LocalizeTemplate(id string, templateData map[string]any) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{
		MessageID:    id,
		TemplateData: templateData,
	})
	if err != nil {
		return notLocalized
	}
	return localization
}

// LocalizeTemplateSingle loads localization using template which has only one key-value pair.
//
// Use only if you called SetLocalizerInstance().
func LocalizeTemplateSingle(id, singleKey string, singleValue any) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{
		MessageID:    id,
		TemplateData: map[string]any{singleKey: singleValue},
	})
	if err != nil {
		return notLocalized
	}
	return localization
}

// LocalizeTemplateSingleWithPlural loads localization using template
// which has only one key-value pair for text with plural count.
//
// Use only if you called SetLocalizerInstance().
func LocalizeTemplateSingleWithPlural(id string, count int, key string, value any) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{
		MessageID: id,
		TemplateData: map[string]any{
			key:           value,
			"PluralCount": count,
		},
		PluralCount: count,
	})
	if err != nil {
		return notLocalized
	}
	return localization
}

// SetABCParams defines the default keys for TemplateData when using an x-func.
// Defining a non-empty slice unlocks the following funcs:
//   - LocalizeTemplateX()
//   - LocalizeTemplateXPlural()
//
// Use only if you called SetLocalizerInstance().
func SetABCParams(abc []string) {
	abcParams = abc
}

func buildABCTemplateData(values []any) map[string]any {
	if len(values) > len(abcParams) {
		return nil
	}

	data := map[string]any{}
	for i, val := range values {
		data[abcParams[i]] = val
	}
	return data
}

// LocalizeTemplateX assigns the given values to the TemplateData keys
// as per the key definitions you set by calling SetABCParams().
//
// Use only if you called SetLocalizerInstance().
func LocalizeTemplateX(id string, values ...any) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{
		MessageID:    id,
		TemplateData: buildABCTemplateData(values),
	})
	if err != nil {
		return notLocalized
	}
	return localization
}

// LocalizeTemplateXPlural assigns the given values to the TemplateData keys
// as per the key definitions you set by calling SetABCParams() with plural count.
//
// Use only if you called SetLocalizerInstance().
func LocalizeTemplateXPlural(id string, count int, values ...any) string {
	localizerMutex.Lock()
	defer localizerMutex.Unlock()

	data := buildABCTemplateData(values)
	data["PluralCount"] = count

	localization, err := GetLocalizerInstance().Localize(&LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		return notLocalized
	}
	return localization
}
