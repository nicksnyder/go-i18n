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

	floatTests := appendFloatTests(nil, 0.1, 0.9, plural.Other)
	floatTests = appendFloatTests(floatTests, 1.1, 1.9, plural.Other)
	floatTests = appendFloatTests(floatTests, 2.1, 2.9, plural.Other)
	floatTests = appendFloatTests(floatTests, 3.1, 3.9, plural.Other)

	language := LanguageWithID("ar")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestCatalan(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 10.0, plural.Other)

	language := LanguageWithID("ca")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestChineseSimplified(t *testing.T) {
	testAlwaysOther(t, LanguageWithID("zh-Hans"))
}

func TestChineseTraditional(t *testing.T) {
	testAlwaysOther(t, LanguageWithID("zh-Hant"))
}

func TestCzech(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Few},
		{3, plural.Few},
		{4, plural.Few},
		{5, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 10.0, plural.Many)

	language := LanguageWithID("cs")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestDanish(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.1, 1.9, plural.One)
	floatTests = appendFloatTests(floatTests, 2.0, 10.0, plural.Other)

	language := LanguageWithID("da")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestEnglish(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 10.0, plural.Other)

	language := LanguageWithID("en")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestFrench(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.One},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 1.9, plural.One)
	floatTests = appendFloatTests(floatTests, 2.0, 10.0, plural.Other)

	language := LanguageWithID("fr")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestGerman(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 10.0, plural.Other)

	language := LanguageWithID("de")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestItalian(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 10.0, plural.Other)

	language := LanguageWithID("it")
	testInts(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestJapanese(t *testing.T) {
	testAlwaysOther(t, LanguageWithID("ja"))
}

func TestSpanish(t *testing.T) {
	intTests := []intPluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}

	floatTests := appendFloatTests(nil, 0.0, 0.9, plural.Other)
	floatTests = appendFloatTests(floatTests, 1.1, 10.0, plural.Other)

	language := LanguageWithID("es")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

// Tests that a language treats all numbers the same.
func testAlwaysOther(t *testing.T, l *Language) {
	testInts(t, l, appendIntTests(nil, 0, 10, plural.Other))
	testFloats(t, l, appendFloatTests(nil, 0.0, 10.0, plural.Other))
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

func appendIntTests(tests []intPluralTest, from, to int64, pc plural.Category) []intPluralTest {
	for i := from; i <= to; i++ {
		tests = append(tests, intPluralTest{i, pc})
	}
	return tests
}

func appendFloatTests(tests []floatPluralTest, from, to float64, pc plural.Category) []floatPluralTest {
	stride := 0.1
	iterations := int64((to - from) / stride)
	for i := int64(0); i < iterations; i++ {
		tests = append(tests, floatPluralTest{from + float64(i)*stride, pc})
	}
	tests = append(tests, floatPluralTest{to, pc})
	return tests
}
