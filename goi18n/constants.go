package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

type constantsCommand struct {
	translationFile string
	packageName     string
	outdir          string
}

const header = `// DON'T CHANGE THIS FILE MANUALLY
// This file was generate using the command:
// $ goi18n constants

`

func (cc *constantsCommand) execute() error {
	if cc.translationFile == "" {
		return fmt.Errorf("need one translation file to parse")
	}

	bundle := bundle.New()

	if err := bundle.LoadTranslationFile(cc.translationFile); err != nil {
		return fmt.Errorf("failed to load translation file %s because %s\n", cc.translationFile, err)
	}

	var buf bytes.Buffer
	buf.WriteString(header)

	buf.WriteString(fmt.Sprintf("package %s\n", cc.packageName))

	var name string
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

	for _, id := range keys {
		name = toCamelCase(id)
		if name == "" {
			return fmt.Errorf("failed on convertion id to constant, id:%s\n", id)
		}

		buf.WriteString(fmt.Sprintf("\n// %s id:%s\n", name, id))
		buf.WriteString(fmt.Sprintf(`const %s string = "%s"%s`, name, id, "\n"))
	}

	filename := filepath.Join(cc.outdir, fmt.Sprintf("%s.go", cc.packageName))
	if err := ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
		return fmt.Errorf("failed to write %s because %s", filename, err)
	}

	return nil
}

// common acronyms used to const names
// https://github.com/golang/lint/blob/master/lint.go
var acronyms = map[string]bool{
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
		if acronyms[w] {
			result += w
		} else {
			result += string(w[0]) + strings.ToLower(string(w[1:]))
		}
	}
	return result
}
