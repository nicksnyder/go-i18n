package main

import "github.com/nicksnyder/go-i18n/v2/i18n"

type MessagList struct {
	HelloPerson        string
	MyUnreadEmails     string
	PersonUnreadEmails string
}

func FillLocalizer(localizer *i18n.Localizer, unreadEmailCount, name string) *MessagList {
	ml := &MessagList{}

	ml.HelloPerson = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    ID1,
			Other: MSG_HELLO,
		},
		TemplateData: map[string]string{
			"Name": name,
		},
	})

	ml.MyUnreadEmails = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          ID2,
			Description: "The number of unread emails I have",
			One:         "I have {{.PluralCount}} unread email.",
			Other:       "I have {{.PluralCount}} unread emails.",
		},
		PluralCount: unreadEmailCount,
	})

	ml.PersonUnreadEmails = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          ID3,
			Description: "The number of unread emails a person has",
			One:         "{{.Name}} has {{.UnreadEmailCount}} unread email.",
			Other:       "{{.Name}} has {{.UnreadEmailCount}} unread emails.",
		},
		PluralCount: unreadEmailCount,
		TemplateData: map[string]interface{}{
			"Name":             name,
			"UnreadEmailCount": unreadEmailCount,
		},
	})

	return ml
}
