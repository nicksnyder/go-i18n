package main

import (
	"flag"
	"fmt"
	"github.com/nicksnyder/go-i18n/src/pkg/i18n"
)

var (
	HelloWorld   = i18n.NewMessage("Hello world!", "This message is displayed when the program begins")
	GoodbyeWorld = i18n.NewMessage("Goodbye world.", "This message is displayed when the program ends")
)

var locale string

func main() {
	flag.StringVar(&locale, "locale", "", "The locale to use for translated messages.")
	flag.Parse()
	fmt.Println(locale)
	i18n.SetLocale(locale)
	fmt.Println(HelloWorld.String())
	fmt.Println(GoodbyeWorld.String())
}
