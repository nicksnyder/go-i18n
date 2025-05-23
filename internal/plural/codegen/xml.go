package main

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// SupplementalData is the top level struct of plural.xml
type SupplementalData struct {
	XMLName      xml.Name      `xml:"supplementalData"`
	PluralGroups []PluralGroup `xml:"plurals>pluralRules"`
}

// PluralGroup is a group of locales with the same plural rules.
type PluralGroup struct {
	Locales     string       `xml:"locales,attr"`
	PluralRules []PluralRule `xml:"pluralRule"`
}

// Name returns a unique name for this plural group.
func (pg *PluralGroup) Name() string {
	n := cases.Title(language.AmericanEnglish).String(pg.Locales)
	return strings.ReplaceAll(n, " ", "")
}

// SplitLocales returns all the locales in the PluralGroup as a slice.
func (pg *PluralGroup) SplitLocales() []string {
	return strings.Split(pg.Locales, " ")
}

// PluralRule is the rule for a single plural form.
type PluralRule struct {
	Count string `xml:"count,attr"`
	Rule  string `xml:",innerxml"`
}

// CountTitle returns the title case of the PluralRule's count.
func (pr *PluralRule) CountTitle() string {
	return cases.Title(language.AmericanEnglish).String(pr.Count)
}

// Condition returns the condition where the PluralRule applies.
func (pr *PluralRule) Condition() string {
	i := strings.Index(pr.Rule, "@")
	return pr.Rule[:i]
}

// Examples returns the integer and decimal examples for the PluralRule.
func (pr *PluralRule) Examples() (integers []string, decimals []string) {
	ex := strings.ReplaceAll(pr.Rule, ", …", "")
	ddelim := "@decimal"
	if i := strings.Index(ex, ddelim); i > 0 {
		dex := strings.TrimSpace(ex[i+len(ddelim):])
		dex = strings.ReplaceAll(dex, "c", "e")
		decimals = strings.Split(dex, ", ")
		ex = ex[:i]
	}
	idelim := "@integer"
	if i := strings.Index(ex, idelim); i > 0 {
		iex := strings.TrimSpace(ex[i+len(idelim):])
		integers = strings.Split(iex, ", ")
		for j, integer := range integers {
			ii := strings.IndexAny(integer, "eEcC")
			if ii > 0 {
				zeros, err := strconv.ParseInt(integer[ii+1:], 10, 0)
				if err != nil {
					panic(err)
				}
				integers[j] = integer[:ii] + strings.Repeat("0", int(zeros))
			}
		}
	}
	return integers, decimals
}

// IntegerExamples returns the integer examples for the PluralRule.
func (pr *PluralRule) IntegerExamples() []string {
	integer, _ := pr.Examples()
	return integer
}

// DecimalExamples returns the decimal examples for the PluralRule.
func (pr *PluralRule) DecimalExamples() []string {
	_, decimal := pr.Examples()
	return decimal
}

var relationRegexp = regexp.MustCompile(`([niftvwce])(?:\s*%\s*([0-9]+))?\s*(!=|=)(.*)`)

// GoCondition converts the XML condition to valid Go code.
func (pr *PluralRule) GoCondition() string {
	var ors []string
	for _, and := range strings.Split(pr.Condition(), "or") {
		var ands []string
		for _, relation := range strings.Split(and, "and") {
			parts := relationRegexp.FindStringSubmatch(relation)
			if parts == nil {
				continue
			}
			lvar := cases.Title(language.AmericanEnglish).String(parts[1])
			lmod, op, rhs := parts[2], parts[3], strings.TrimSpace(parts[4])
			if op == "=" {
				op = "=="
			}
			if lvar == "E" {
				// E is a deprecated symbol for C
				// https://unicode.org/reports/tr35/tr35-numbers.html#Plural_Operand_Meanings
				lvar = "C"
			}
			lvar = "ops." + lvar
			var rhor []string
			var rany []string
			for _, rh := range strings.Split(rhs, ",") {
				if parts := strings.Split(rh, ".."); len(parts) == 2 {
					from, to := parts[0], parts[1]
					if lvar == "ops.N" {
						if lmod != "" {
							rhor = append(rhor, fmt.Sprintf("ops.NModInRange(%s, %s, %s)", lmod, from, to))
						} else {
							rhor = append(rhor, fmt.Sprintf("ops.NInRange(%s, %s)", from, to))
						}
					} else if lmod != "" {
						rhor = append(rhor, fmt.Sprintf("intInRange(%s %% %s, %s, %s)", lvar, lmod, from, to))
					} else {
						rhor = append(rhor, fmt.Sprintf("intInRange(%s, %s, %s)", lvar, from, to))
					}
				} else {
					rany = append(rany, rh)
				}
			}

			if len(rany) > 0 {
				rh := strings.Join(rany, ",")
				if lvar == "ops.N" {
					if lmod != "" {
						rhor = append(rhor, fmt.Sprintf("ops.NModEqualsAny(%s, %s)", lmod, rh))
					} else {
						rhor = append(rhor, fmt.Sprintf("ops.NEqualsAny(%s)", rh))
					}
				} else if lmod != "" {
					rhor = append(rhor, fmt.Sprintf("intEqualsAny(%s %% %s, %s)", lvar, lmod, rh))
				} else {
					rhor = append(rhor, fmt.Sprintf("intEqualsAny(%s, %s)", lvar, rh))
				}
			}
			r := strings.Join(rhor, " || ")
			if len(rhor) > 1 {
				r = "(" + r + ")"
			}
			if op == "!=" {
				r = "!" + r
			}
			ands = append(ands, r)
		}
		ors = append(ors, strings.Join(ands, " && "))
	}
	return strings.Join(ors, " ||\n")
}
