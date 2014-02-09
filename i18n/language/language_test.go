package language

import (
	"github.com/nicksnyder/go-i18n/i18n/plural"
	"testing"
)

type intPluralTest struct {
	i  int64
	pc plural.Category
}

type floatPluralTest struct {
	f  float64
	pc plural.Category
}

func TestArabic(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Zero},
		{1, plural.One},
		{2, plural.Two},
		{3, plural.Few},
		{10, plural.Few},
		{103, plural.Few},
		{110, plural.Few},
		{11, plural.Many},
		{99, plural.Many},
		{111, plural.Many},
		{199, plural.Many},
		{100, plural.Other},
		{102, plural.Other},
		{200, plural.Other},
		{202, plural.Other},
	}

	floatTests := []floatPluralTest{
		{0.1, plural.Other},
		{0.2, plural.Other},
		{0.3, plural.Other},
		{1.1, plural.Other},
		{1.2, plural.Other},
		{1.3, plural.Other},
		{2.1, plural.Other},
		{2.2, plural.Other},
		{2.3, plural.Other},
		{3.1, plural.Other},
		{3.2, plural.Other},
		{3.3, plural.Other},
	}

	language := LanguageWithID("ar")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestCatalan(t *testing.T) {
	testOneIsSpecial(t, LanguageWithID("ca"))
}

func TestChineseSimplified(t *testing.T) {
	testAlwaysOther(t, LanguageWithID("zh-Hans"))
}

func TestChineseTraditional(t *testing.T) {
	testAlwaysOther(t, LanguageWithID("zh-Hant"))
}

func TestEnglish(t *testing.T) {
	testOneIsSpecial(t, LanguageWithID("en"))
}

func TestFrench(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.One},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := []floatPluralTest{
		{0.1, plural.One},
		{0.2, plural.One},
		{0.9, plural.One},
		{1.1, plural.One},
		{1.2, plural.One},
		{1.9, plural.One},
		{2.1, plural.Other},
		{2.2, plural.Other},
	}

	language := LanguageWithID("fr")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestGerman(t *testing.T) {
	testOneIsSpecial(t, LanguageWithID("de"))
}

func TestItalian(t *testing.T) {
	testOneIsSpecial(t, LanguageWithID("it"))
}

func TestJapanese(t *testing.T) {
	testAlwaysOther(t, LanguageWithID("ja"))
}

func TestSpanish(t *testing.T) {
	testOneIsSpecial(t, LanguageWithID("es"))
}

// Tests that a language treats one as special and all other numbers the same.
func testOneIsSpecial(t *testing.T, l *Language) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := []floatPluralTest{
		{0.1, plural.Other},
		{0.2, plural.Other},
		{1.1, plural.Other},
		{1.2, plural.Other},
		{2.1, plural.Other},
		{2.2, plural.Other},
	}

	testInts(t, l, intTests)
	testIntsAsFloats(t, l, intTests)
	testFloats(t, l, floatTests)
}

// Tests that a language treats all numbers the same.
func testAlwaysOther(t *testing.T, l *Language) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.Other},
		{2, plural.Other},
	}

	floatTests := []floatPluralTest{
		{0.1, plural.Other},
		{0.2, plural.Other},
		{1.1, plural.Other},
		{1.2, plural.Other},
		{2.1, plural.Other},
		{2.2, plural.Other},
	}

	testInts(t, l, intTests)
	testIntsAsFloats(t, l, intTests)
	testFloats(t, l, floatTests)
}

func testInts(t *testing.T, language *Language, intTests []intPluralTest) {
	for _, test := range intTests {
		if pc := language.int64PluralCategory(test.i); pc != test.pc {
			t.Errorf("Int64PluralCategory(%d) returned %s; expected %s", test.i, pc, test.pc)
		}
	}
}

func testIntsAsFloats(t *testing.T, language *Language, intTests []intPluralTest) {
	for _, test := range intTests {
		f := float64(test.i)
		if pc := language.float64PluralCategory(f); pc != test.pc {
			t.Errorf("Float64PluralCategory(%f) returned %s; expected %s", f, pc, test.pc)
		}
	}
}

func testFloats(t *testing.T, language *Language, floatTests []floatPluralTest) {
	for _, test := range floatTests {
		if pc := language.float64PluralCategory(test.f); pc != test.pc {
			t.Errorf("Float64PluralCategory(%f) returned %s; expected %s", test.f, pc, test.pc)
		}
	}
}
