package main

import (
	"os"
	"fmt"
	"embed"
	"golang.org/x/text/language"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/Xuanwo/go-locale"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

// This will cause go to embed the translation files into the binary when it is compiled.
//go:embed active.*.json
var LocaleFS embed.FS

// Init the localizer
var localizer = InitLocalizer()

// Sets the localizer with the proper language
func InitLocalizer() *goi18n.Localizer {
	locale := DetectLocal()
	bundle := goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	lang_file := fmt.Sprintf("active.%s.json", locale)
	_, err := bundle.LoadMessageFileFS(LocaleFS, lang_file)
	if err != nil {
		fmt.Println("Unable to load language from %s", lang_file)
	}

	loc := goi18n.NewLocalizer(bundle, locale)
	return loc
}

// Tries to determine the system locale, when local isn't set, default to en_US
// Mostly this reads from 'LANGUAGE' ENV variable. https://github.com/Xuanwo/go-locale for other ways to set this.
func DetectLocal() string {
	tag, err := locale.Detect()
	if err != nil {
		return "en-US"
	}
	return tag.String()
}

// Translates a string, with any substitutions needed. Handles simple cases where plural count is not needed.
// text: string to be translated
// subs: A single map[string]interface{} incase the text string needs variable substitutions.
func T(text string, subs ...interface{}) string {

	config := &goi18n.LocalizeConfig{
		DefaultMessage: &goi18n.Message{ID: text, Other: text},
	}
	// Need to use `subs ...interface{}` so that we can have 0 or 1 subs.
	if subs != nil && len(subs) == 1 {
		config.TemplateData = subs[0]
	}

	// Actually translate the message.
	l_string, err := localizer.Localize(config)

	// If a string is not in the translation file, print an error.
	if err != nil {
		// Can't recursively use the T() function here otherwise you get an infinite loop.
		error_message := &goi18n.LocalizeConfig{
			DefaultMessage: &goi18n.Message{
				ID: "T_Error_Message",
				Other: "Translation Error: {{.Error}}",
			},
		}
		error_message.TemplateData = map[string]interface{}{
			"Error": err.Error(),
		}
		e_string, _ := localizer.Localize(error_message)
		fmt.Println(e_string)
	}
	return l_string
}

var rootCmd = &cobra.Command{
	Use: "cli",
	Short: T("This is an example of how to translate a CLI application"),
	Long: T("This is a longer description."),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(T("Hello World!"))
	},
}

func main() {
	// Just an example of how string substitution can work with T() functions
	user_locale := DetectLocal()
	locale_sub := map[string]interface{}{"User_locale": user_locale}
	fmt.Println(T("Detected Locale is {{.User_locale}}", locale_sub))

	// Actually run the command and handle any errors.
	if err := rootCmd.Execute(); err != nil {
		subs := map[string]interface{}{"Error": err.Error()}
		fmt.Println(T("ERROR: {{.Error}}", subs), err.Error())
		os.Exit(1)
	}
}
