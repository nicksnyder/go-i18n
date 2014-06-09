package i18n_test

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/i18n"
)

func Example() {
	i18n.MustLoadTranslationFile("../goi18n/testdata/expected/en-us.all.json")

	T, _ := i18n.Tfunc("en-US")

	fmt.Println(T("program_greeting"))
	fmt.Println(T("person_greeting", map[string]interface{}{
		"Person": "Bob",
	}))

	fmt.Println(T("your_unread_email_count", 0))
	fmt.Println(T("your_unread_email_count", 1))
	fmt.Println(T("your_unread_email_count", 2))
	fmt.Println(T("my_height_in_meters", "1.7"))

	fmt.Println(T("person_unread_email_count", 0, map[string]interface{}{
		"Person": "Bob",
	}))
	fmt.Println(T("person_unread_email_count", 1, map[string]interface{}{
		"Person": "Bob",
	}))
	fmt.Println(T("person_unread_email_count", 2, map[string]interface{}{
		"Person": "Bob",
	}))

	fmt.Println(T("person_unread_email_count_timeframe", 3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": T("d_days", 0),
	}))
	fmt.Println(T("person_unread_email_count_timeframe", 3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": T("d_days", 1),
	}))
	fmt.Println(T("person_unread_email_count_timeframe", 3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": T("d_days", 2),
	}))

	// Output:
	// Hello world
	// Hello Bob
	// You have 0 unread emails.
	// You have 1 unread email.
	// You have 2 unread emails.
	// I am 1.7 meters tall.
	// Bob has 0 unread emails.
	// Bob has 1 unread email.
	// Bob has 2 unread emails.
	// Bob has 3 unread emails in the past 0 days.
	// Bob has 3 unread emails in the past 1 day.
	// Bob has 3 unread emails in the past 2 days.
}
