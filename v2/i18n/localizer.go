package i18n

import (
	"fmt"
	"text/template"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
	"golang.org/x/text/language"
)

// Localizer provides Localize and MustLocalize methods that return localized messages.
// Localize and MustLocalize methods use a language.Tag matching algorithm based
// on the best possible value. This algorithm may cause an unexpected language.Tag returned
// value depending on the order of the tags stored in memory. For example, if the bundle
// used to create a Localizer instance ingested locales following this order
// ["en-US", "en-GB", "en-IE", "en"] and the locale "en" is asked, the underlying matching
// algorithm will return "en-US" thinking it is the best match possible. More information
// about the algorithm in this Github issue: https://github.com/golang/go/issues/49176.
// There is additionnal informations inside the Go code base:
// https://github.com/golang/text/blob/master/language/match.go#L142
type Localizer struct {
	// bundle contains the messages that can be returned by the Localizer.
	bundle *Bundle

	// tags is the list of language tags that the Localizer checks
	// in order when localizing a message.
	tags []language.Tag
}

// NewLocalizer returns a new Localizer that looks up messages
// in the bundle according to the language preferences in langs.
// It can parse Accept-Language headers as defined in http://www.ietf.org/rfc/rfc2616.txt.
func NewLocalizer(bundle *Bundle, langs ...string) *Localizer {
	return &Localizer{
		bundle: bundle,
		tags:   parseTags(langs),
	}
}

func parseTags(langs []string) []language.Tag {
	tags := []language.Tag{}
	for _, lang := range langs {
		t, _, err := language.ParseAcceptLanguage(lang)
		if err != nil {
			continue
		}
		tags = append(tags, t...)
	}
	return tags
}

// LocalizeConfig configures a call to the Localize method on Localizer.
type LocalizeConfig struct {
	// MessageID is the id of the message to lookup.
	// This field is ignored if DefaultMessage is set.
	MessageID string

	// TemplateData is the data passed when executing the message's template.
	// If TemplateData is nil and PluralCount is not nil, then the message template
	// will be executed with data that contains the plural count.
	TemplateData interface{}

	// PluralCount determines which plural form of the message is used.
	PluralCount interface{}

	// DefaultMessage is used if the message is not found in any message files.
	DefaultMessage *Message

	// Funcs is used to extend the Go template engine's built in functions
	Funcs template.FuncMap
}

type invalidPluralCountErr struct {
	messageID   string
	pluralCount interface{}
	err         error
}

func (e *invalidPluralCountErr) Error() string {
	return fmt.Sprintf("invalid plural count %#v for message id %q: %s", e.pluralCount, e.messageID, e.err)
}

// MessageNotFoundErr is returned from Localize when a message could not be found.
type MessageNotFoundErr struct {
	tag       language.Tag
	messageID string
}

func (e *MessageNotFoundErr) Error() string {
	return fmt.Sprintf("message %q not found in language %q", e.messageID, e.tag)
}

type pluralizeErr struct {
	messageID string
	tag       language.Tag
}

func (e *pluralizeErr) Error() string {
	return fmt.Sprintf("unable to pluralize %q because there no plural rule for %q", e.messageID, e.tag)
}

type messageIDMismatchErr struct {
	messageID        string
	defaultMessageID string
}

func (e *messageIDMismatchErr) Error() string {
	return fmt.Sprintf("message id %q does not match default message id %q", e.messageID, e.defaultMessageID)
}

// Localize returns a localized message.
func (l *Localizer) Localize(lc *LocalizeConfig) (string, error) {
	msg, _, err := l.LocalizeWithTag(lc)
	return msg, err
}

// Localize returns a localized message.
func (l *Localizer) LocalizeMessage(msg *Message) (string, error) {
	return l.Localize(&LocalizeConfig{
		DefaultMessage: msg,
	})
}

// TODO: uncomment this (and the test) when extract has been updated to extract these call sites too.
// Localize returns a localized message.
// func (l *Localizer) LocalizeMessageID(messageID string) (string, error) {
// 	return l.Localize(&LocalizeConfig{
// 		MessageID: messageID,
// 	})
// }

// LocalizeWithTag returns a localized message and the language tag.
// It may return a best effort localized message even if an error happens.
func (l *Localizer) LocalizeWithTag(lc *LocalizeConfig) (string, language.Tag, error) {
	messageID := lc.MessageID
	if lc.DefaultMessage != nil {
		if messageID != "" && messageID != lc.DefaultMessage.ID {
			return "", language.Und, &messageIDMismatchErr{messageID: messageID, defaultMessageID: lc.DefaultMessage.ID}
		}
		messageID = lc.DefaultMessage.ID
	}

	var operands *plural.Operands
	templateData := lc.TemplateData
	if lc.PluralCount != nil {
		var err error
		operands, err = plural.NewOperands(lc.PluralCount)
		if err != nil {
			return "", language.Und, &invalidPluralCountErr{messageID: messageID, pluralCount: lc.PluralCount, err: err}
		}
		if templateData == nil {
			templateData = map[string]interface{}{
				"PluralCount": lc.PluralCount,
			}
		}
	}

	tag, template, err := l.getMessageTemplate(messageID, lc.DefaultMessage)
	if template == nil {
		return "", language.Und, err
	}

	pluralForm := l.pluralForm(tag, operands)
	msg, err2 := template.Execute(pluralForm, templateData, lc.Funcs)
	if err2 != nil {
		if err == nil {
			err = err2
		}

		// Attempt to fallback to "Other" pluralization in case translations are incomplete.
		if pluralForm != plural.Other {
			msg2, err3 := template.Execute(plural.Other, templateData, lc.Funcs)
			if err3 == nil {
				msg = msg2
			}
		}
	}
	return msg, tag, err
}

func (l *Localizer) getMessageTemplate(id string, defaultMessage *Message) (language.Tag, *MessageTemplate, error) {
	_, i, _ := l.bundle.matcher.Match(l.tags...)
	tag := l.bundle.tags[i]
	mt := l.bundle.getMessageTemplate(tag, id)
	if mt != nil {
		return tag, mt, nil
	}

	if tag == l.bundle.defaultLanguage {
		if defaultMessage == nil {
			return language.Und, nil, &MessageNotFoundErr{tag: tag, messageID: id}
		}
		mt := NewMessageTemplate(defaultMessage)
		if mt == nil {
			return language.Und, nil, &MessageNotFoundErr{tag: tag, messageID: id}
		}
		return tag, mt, nil
	}

	// Fallback to default language in bundle.
	mt = l.bundle.getMessageTemplate(l.bundle.defaultLanguage, id)
	if mt != nil {
		return l.bundle.defaultLanguage, mt, &MessageNotFoundErr{tag: tag, messageID: id}
	}

	// Fallback to default message.
	if defaultMessage == nil {
		return language.Und, nil, &MessageNotFoundErr{tag: tag, messageID: id}
	}
	return l.bundle.defaultLanguage, NewMessageTemplate(defaultMessage), &MessageNotFoundErr{tag: tag, messageID: id}
}

func (l *Localizer) pluralForm(tag language.Tag, operands *plural.Operands) plural.Form {
	if operands == nil {
		return plural.Other
	}
	return l.bundle.pluralRules.Rule(tag).PluralFormFunc(operands)
}

// MustLocalize is similar to Localize, except it panics if an error happens.
func (l *Localizer) MustLocalize(lc *LocalizeConfig) string {
	localized, err := l.Localize(lc)
	if err != nil {
		panic(err)
	}
	return localized
}
