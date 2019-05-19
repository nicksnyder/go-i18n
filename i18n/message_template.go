package i18n

import (
	"bytes"
	"fmt"

	"text/template"

	"github.com/nicksnyder/go-i18n/v2/internal"
	"github.com/nicksnyder/go-i18n/v2/internal/plural"
)

// MessageTemplate is an executable template for a message.
type MessageTemplate struct {
	*Message
	PluralTemplates map[plural.Form]*internal.Template
}

// NewMessageTemplate returns a new message template.
func NewMessageTemplate(m *Message) *MessageTemplate {
	pluralTemplates := map[plural.Form]*internal.Template{}
	setPluralTemplate(pluralTemplates, plural.Zero, m.Zero)
	setPluralTemplate(pluralTemplates, plural.One, m.One)
	setPluralTemplate(pluralTemplates, plural.Two, m.Two)
	setPluralTemplate(pluralTemplates, plural.Few, m.Few)
	setPluralTemplate(pluralTemplates, plural.Many, m.Many)
	setPluralTemplate(pluralTemplates, plural.Other, m.Other)
	if len(pluralTemplates) == 0 {
		return nil
	}
	return &MessageTemplate{
		Message:         m,
		PluralTemplates: pluralTemplates,
	}
}

func setPluralTemplate(pluralTemplates map[plural.Form]*internal.Template, pluralForm plural.Form, src string) {
	if src != "" {
		pluralTemplates[pluralForm] = &internal.Template{Src: src}
	}
}

type pluralFormNotFoundError struct {
	pluralForm plural.Form
	messageID  string
}

func (e pluralFormNotFoundError) Error() string {
	return fmt.Sprintf("message %q has no plural form %q", e.messageID, e.pluralForm)
}

// Execute executes the template for the plural form and template data.
func (mt *MessageTemplate) Execute(pluralForm plural.Form, data interface{}, funcs template.FuncMap) (string, error) {
	t := mt.PluralTemplates[pluralForm]
	if t == nil {
		return "", pluralFormNotFoundError{
			pluralForm: pluralForm,
			messageID:  mt.Message.ID,
		}
	}
	if err := t.Parse(mt.LeftDelim, mt.RightDelim, funcs); err != nil {
		return "", err
	}
	if t.Template == nil {
		return t.Src, nil
	}
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
