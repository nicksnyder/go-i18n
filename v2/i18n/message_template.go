package i18n

import "bytes"

// MessageTemplate is an executable template for a message.
type MessageTemplate struct {
	*Message
	PluralTemplates map[PluralForm]*Template
}

// NewMessageTemplate returns a new message template.
func NewMessageTemplate(m *Message) *MessageTemplate {
	mt := &MessageTemplate{
		Message:         m,
		PluralTemplates: make(map[PluralForm]*Template),
	}
	mt.setPluralTemplate(Zero, m.Zero)
	mt.setPluralTemplate(One, m.One)
	mt.setPluralTemplate(Two, m.Two)
	mt.setPluralTemplate(Few, m.Few)
	mt.setPluralTemplate(Many, m.Many)
	mt.setPluralTemplate(Other, m.Other)
	return mt
}

func (mt *MessageTemplate) setPluralTemplate(pluralForm PluralForm, src string) {
	if src != "" {
		mt.PluralTemplates[pluralForm] = &Template{Src: src}
	}
}

// Execute executes the template for the plural form and template data.
func (mt *MessageTemplate) Execute(pluralForm PluralForm, data interface{}) (string, error) {
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
