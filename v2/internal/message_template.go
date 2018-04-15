package internal

import (
	"bytes"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
)

// MessageTemplate is an executable template for a message.
type MessageTemplate struct {
	*Message
	PluralTemplates map[plural.Form]*Template
}

// NewMessageTemplate returns a new message template.
func NewMessageTemplate(m *Message) *MessageTemplate {
	mt := &MessageTemplate{
		Message:         m,
		PluralTemplates: make(map[plural.Form]*Template),
	}
	mt.setPluralTemplate(plural.Zero, m.Zero)
	mt.setPluralTemplate(plural.One, m.One)
	mt.setPluralTemplate(plural.Two, m.Two)
	mt.setPluralTemplate(plural.Few, m.Few)
	mt.setPluralTemplate(plural.Many, m.Many)
	mt.setPluralTemplate(plural.Other, m.Other)
	return mt
}

func (mt *MessageTemplate) setPluralTemplate(pluralForm plural.Form, src string) {
	if src != "" {
		mt.PluralTemplates[pluralForm] = &Template{Src: src}
	}
}

// Execute executes the template for the plural form and template data.
func (mt *MessageTemplate) Execute(pluralForm plural.Form, data interface{}) (string, error) {
	t := mt.PluralTemplates[pluralForm]
	if err := t.parse(mt.LeftDelim, mt.RightDelim); err != nil {
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
