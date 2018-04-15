// Command example runs a sample webserver that uses go-i18n/v2/i18n.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var page = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<body>

<h1>{{.Title}}</h1>

{{range .Paragraphs}}<p>{{.}}</p>{{end}}

</body>
</html>
`))

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		name := r.FormValue("name")
		if name == "" {
			name = "Bob"
		}

		unreadEmailCount, _ := strconv.ParseInt(r.FormValue("unreadEmailCount"), 10, 64)
		catCount, _ := strconv.ParseInt(r.FormValue("catCount"), 10, 64)
		catColor := r.FormValue("catColor")
		if catColor == "" {
			catColor = "black"
		}

		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "HelloPerson",
				Other: "Hello {{.Name}}",
			},
			TemplateData: map[string]string{
				"Name": name,
			},
		})

		unreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "UnreadEmails",
				Description: "The number of unread emails a person has",
				One:         "I have {{.PluralCount}} unread email.",
				Other:       "I have {{.PluralCount}} unread emails.",
			},
			PluralCount: unreadEmailCount,
		})

		coloredCats := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "ColoredCats",
				Description: "The number of cats a person has",
				One:         "I have {{.Count}} {{.Color}} cat.",
				Other:       "I have {{.Count}} {{.Color}} cats.",
			},
			PluralCount: catCount,
			TemplateData: map[string]interface{}{
				"Count": catCount,
				"Color": catColor,
			},
		})

		err := page.Execute(w, map[string]interface{}{
			"Title": helloPerson,
			"Paragraphs": []string{
				unreadEmails,
				coloredCats,
			},
		})
		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
