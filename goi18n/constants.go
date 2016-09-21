package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

type constantsCommand struct {
	translationFiles []string
	packageName      string
	outdir           string
}

type templateConstants struct {
	ID   string
	Name string
}

type templateHeader struct {
	PackageName string
	Constants   []templateConstants
}

var constTemplate = template.Must(template.New("").Parse(`// DON'T CHANGE THIS FILE MANUALLY
// This file was generated using the command:
// $ goi18n constants

package {{.PackageName}}
{{range .Constants}}
const {{.Name}} = "{{.ID}}"
{{end}}`))

func (cc *constantsCommand) execute() error {
	if len(cc.translationFiles) != 1 {
		return fmt.Errorf("need one translation file")
	}

	bundle := bundle.New()

	if err := bundle.LoadTranslationFile(cc.translationFiles[0]); err != nil {
		return fmt.Errorf("failed to load translation file %s because %s\n", cc.translationFiles[0], err)
	}

	translations := bundle.Translations()
	lang := translations[bundle.LanguageTags()[0]]

	// create an array of id to organize
	keys := make([]string, len(lang))
	i := 0

	for id := range lang {
		keys[i] = id
		i++
	}
	sort.Strings(keys)

	tmpl := &templateHeader{
		PackageName: cc.packageName,
		Constants:   make([]templateConstants, len(keys)),
	}

	for i, id := range keys {
		tmpl.Constants[i].ID = id
		tmpl.Constants[i].Name = toCamelCase(id)
	}

	filename := filepath.Join(cc.outdir, fmt.Sprintf("%s.go", cc.packageName))
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s because %s", filename, err)
	}

	defer f.Close()

	if err = constTemplate.Execute(f, tmpl); err != nil {
		return fmt.Errorf("failed to write file %s because %s", filename, err)
	}

	return nil
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
// https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

func toCamelCase(id string) string {
	r := regexp.MustCompile("[A-Za-z0-9]+")
	words := r.FindAllString(id, -1)
	var result string
	for _, w := range words {
		w = strings.ToUpper(w)
		if commonInitialisms[w] {
			result += w
		} else {
			result += string(w[0]) + strings.ToLower(string(w[1:]))
		}
	}
	return result
}
