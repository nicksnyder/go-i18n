package i18n

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
	"golang.org/x/text/language"
)

type localizerTest struct {
	name              string
	defaultLanguage   language.Tag
	messages          map[language.Tag][]*Message
	acceptLangs       []string
	conf              *LocalizeConfig
	expectedErr       error
	expectedLocalized string
}

func localizerTests() []localizerTest {
	return []localizerTest{
		{
			name:            "message id mismatch",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloWorld",
				DefaultMessage: &Message{
					ID: "DefaultHelloWorld",
				},
			},
			expectedErr: &messageIDMismatchErr{messageID: "HelloWorld", defaultMessageID: "DefaultHelloWorld"},
		},
		{
			name:            "message id not mismatched",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "Hello!"}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloWorld",
				DefaultMessage: &Message{
					ID: "HelloWorld",
				},
			},
			expectedLocalized: "Hello!",
		},
		{
			name:              "missing translation from default language",
			defaultLanguage:   language.English,
			acceptLangs:       []string{"en"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{tag: language.English, messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "empty translation without fallback",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "Hello World!"}},
				language.Spanish: {{ID: "HelloWorld"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{tag: language.Spanish, messageID: "HelloWorld"},
			expectedLocalized: "Hello World!",
		},
		{
			name:            "missing translation from default language with other translation",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"en"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{tag: language.English, messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:              "missing translations from not default language",
			defaultLanguage:   language.English,
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{tag: language.English, messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "missing translation from not default language",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "SomethingElse", Other: "other"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{tag: language.Spanish, messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "missing translation not default language with other translation",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.French:  {{ID: "HelloWorld", Other: "other"}},
				language.Spanish: {{ID: "SomethingElse", Other: "other"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedErr:       &MessageNotFoundErr{tag: language.Spanish, messageID: "HelloWorld"},
			expectedLocalized: "",
		},
		{
			name:            "accept default language, message in bundle",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"en"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedLocalized: "other",
		},
		{
			name:            "accept default language, message in bundle, default message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "accept not default language, message in bundle",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "other"}},
			},
			acceptLangs:       []string{"es"},
			conf:              &LocalizeConfig{MessageID: "HelloWorld"},
			expectedLocalized: "other",
		},
		{
			name:            "accept not default language, other message in bundle, default message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "accept not default language, message in bundle, default message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "accept default language, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "default other",
		},
		{
			name:            "accept not default language, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "default other",
		},
		{
			name:            "fallback to non-default less specific language",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.Spanish: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es-ES"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "fallback to non-default more specific language",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.EuropeanSpanish: {{ID: "HelloWorld", Other: "bundle other"}},
			},
			acceptLangs: []string{"es"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{ID: "HelloWorld", Other: "default other"},
			},
			expectedLocalized: "bundle other",
		},
		{
			name:            "plural count one, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID:   "Cats",
				PluralCount: 1,
			},
			expectedLocalized: "I have 1 cat",
		},
		{
			name:            "plural count other, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID:   "Cats",
				PluralCount: 2,
			},
			expectedLocalized: "I have 2 cats",
		},
		{
			name:            "plural count float, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID:   "Cats",
				PluralCount: "2.5",
			},
			expectedLocalized: "I have 2.5 cats",
		},
		{
			name:            "plural count one, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: 1,
				DefaultMessage: &Message{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 1 cat",
		},
		{
			name:            "plural count missing one, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: 1,
				DefaultMessage: &Message{
					ID:    "Cats",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 1 cats",
			expectedErr:       pluralFormNotFoundError{messageID: "Cats", pluralForm: plural.One},
		},
		{
			name:            "plural count missing other, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: 2,
				DefaultMessage: &Message{
					ID:  "Cats",
					One: "I have {{.PluralCount}} cat",
				},
			},
			expectedLocalized: "",
			expectedErr:       pluralFormNotFoundError{messageID: "Cats", pluralForm: plural.Other},
		},
		{
			name:            "plural count other, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: 2,
				DefaultMessage: &Message{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 2 cats",
		},
		{
			name:            "plural count float, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				PluralCount: "2.5",
				DefaultMessage: &Message{
					ID:    "Cats",
					One:   "I have {{.PluralCount}} cat",
					Other: "I have {{.PluralCount}} cats",
				},
			},
			expectedLocalized: "I have 2.5 cats",
		},
		{
			name:            "template data, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "HelloPerson",
					Other: "Hello {{.Person}}",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloPerson",
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "HelloPerson",
					Other: "Hello {{.Person}}",
				},
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, custom delims, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:         "HelloPerson",
					Other:      "Hello <<.Person>>",
					LeftDelim:  "<<",
					RightDelim: ">>",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "HelloPerson",
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, custom delims, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:         "HelloPerson",
					Other:      "Hello <<.Person>>",
					LeftDelim:  "<<",
					RightDelim: ">>",
				},
				TemplateData: map[string]string{
					"Person": "Nick",
				},
			},
			expectedLocalized: "Hello Nick",
		},
		{
			name:            "template data, plural count one, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "PersonCats",
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  1,
				},
				PluralCount: 1,
			},
			expectedLocalized: "Nick has 1 cat",
		},
		{
			name:            "template data, plural count other, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "PersonCats",
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  2,
				},
				PluralCount: 2,
			},
			expectedLocalized: "Nick has 2 cats",
		},
		{
			name:            "template data, plural count float, bundle message",
			defaultLanguage: language.English,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				}},
			},
			acceptLangs: []string{"en"},
			conf: &LocalizeConfig{
				MessageID: "PersonCats",
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  "2.5",
				},
				PluralCount: "2.5",
			},
			expectedLocalized: "Nick has 2.5 cats",
		},
		{
			name:            "template data, plural count one, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				},
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  1,
				},
				PluralCount: 1,
			},
			expectedLocalized: "Nick has 1 cat",
		},
		{
			name:            "template data, plural count other, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				},
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  2,
				},
				PluralCount: 2,
			},
			expectedLocalized: "Nick has 2 cats",
		},
		{
			name:            "template data, plural count float, default message",
			defaultLanguage: language.English,
			acceptLangs:     []string{"en"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "PersonCats",
					One:   "{{.Person}} has {{.Count}} cat",
					Other: "{{.Person}} has {{.Count}} cats",
				},
				TemplateData: map[string]interface{}{
					"Person": "Nick",
					"Count":  "2.5",
				},
				PluralCount: "2.5",
			},
			expectedLocalized: "Nick has 2.5 cats",
		},
		{
			name:            "no fallback",
			defaultLanguage: language.Spanish,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Hello",
					Other: "Hello!",
				}},
				language.AmericanEnglish: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
			},
			acceptLangs: []string{"en-US"},
			conf: &LocalizeConfig{
				MessageID: "Hello",
			},
			expectedErr: &MessageNotFoundErr{tag: language.AmericanEnglish, messageID: "Hello"},
		},
		{
			name:            "fallback default message",
			defaultLanguage: language.Spanish,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
				language.AmericanEnglish: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
			},
			acceptLangs: []string{"en-US"},
			conf: &LocalizeConfig{
				DefaultMessage: &Message{
					ID:    "Hello",
					Other: "Hola!",
				},
			},
			expectedLocalized: "Hola!",
			expectedErr:       &MessageNotFoundErr{tag: language.AmericanEnglish, messageID: "Hello"},
		},
		{
			name:            "no fallback default message",
			defaultLanguage: language.Spanish,
			messages: map[language.Tag][]*Message{
				language.English: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
				language.AmericanEnglish: {{
					ID:    "Goodbye",
					Other: "Goodbye!",
				}},
			},
			acceptLangs: []string{"en-US"},
			conf: &LocalizeConfig{
				MessageID: "Hello",
			},
			expectedErr: &MessageNotFoundErr{tag: language.AmericanEnglish, messageID: "Hello"},
		},
	}
}

func TestLocalizer_Localize(t *testing.T) {
	for _, test := range localizerTests() {
		t.Run(test.name, func(t *testing.T) {
			bundle := NewBundle(test.defaultLanguage)
			for tag, messages := range test.messages {
				if err := bundle.AddMessages(tag, messages...); err != nil {
					t.Fatal(err)
				}
			}
			check := func(localized string, err error) {
				t.Helper()
				if !reflect.DeepEqual(err, test.expectedErr) {
					t.Errorf("expected error %#v; got %#v", test.expectedErr, err)
				}
				if localized != test.expectedLocalized {
					t.Errorf("expected localized string %q; got %q", test.expectedLocalized, localized)
				}
			}
			localizer := NewLocalizer(bundle, test.acceptLangs...)
			check(localizer.Localize(test.conf))

			if test.conf.DefaultMessage != nil && reflect.DeepEqual(test.conf, &LocalizeConfig{DefaultMessage: test.conf.DefaultMessage}) {
				check(localizer.LocalizeMessage(test.conf.DefaultMessage))
			}

			// if test.conf.MessageID != "" && reflect.DeepEqual(test.conf, &LocalizeConfig{MessageID: test.conf.MessageID}) {
			// 	check(localizer.LocalizeMessageID(test.conf.MessageID))
			// }
		})
	}
}

func BenchmarkLocalizer_Localize(b *testing.B) {
	for _, test := range localizerTests() {
		b.Run(test.name, func(b *testing.B) {
			bundle := NewBundle(test.defaultLanguage)
			for tag, messages := range test.messages {
				if err := bundle.AddMessages(tag, messages...); err != nil {
					b.Fatal(err)
				}
			}

			localizer := NewLocalizer(bundle, test.acceptLangs...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = localizer.Localize(test.conf)
			}
		})
	}
}

func TestMessageNotFoundError(t *testing.T) {
	actual := (&MessageNotFoundErr{tag: language.AmericanEnglish, messageID: "hello"}).Error()
	expected := `message "hello" not found in language "en-US"`
	if actual != expected {
		t.Fatalf("expected %q; got %q", expected, actual)
	}
}

func TestMessageIDMismatchError(t *testing.T) {
	actual := (&messageIDMismatchErr{messageID: "hello", defaultMessageID: "world"}).Error()
	expected := `message id "hello" does not match default message id "world"`
	if actual != expected {
		t.Fatalf("expected %q; got %q", expected, actual)
	}
}

func TestInvalidPluralCountError(t *testing.T) {
	actual := (&invalidPluralCountErr{messageID: "hello", pluralCount: "blah", err: fmt.Errorf("error")}).Error()
	expected := `invalid plural count "blah" for message id "hello": error`
	if actual != expected {
		t.Fatalf("expected %q; got %q", expected, actual)
	}
}

func TestMustLocalize(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustLocalize did not panic")
		}
	}()
	bundle := NewBundle(language.English)
	localizer := NewLocalizer(bundle)
	localizer.MustLocalize(&LocalizeConfig{
		MessageID: "hello",
	})
}
