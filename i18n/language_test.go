package i18n

import (
	"testing"
)

type intPluralTest struct {
	i  int64
	pc PluralCategory
}

type floatPluralTest struct {
	f  float64
	pc PluralCategory
}

func TestArabic(t *testing.T) {
	intTests := []intPluralTest{
		{0, Zero},
		{1, One},
		{2, Two},
		{3, Few},
		{10, Few},
		{103, Few},
		{110, Few},
		{11, Many},
		{99, Many},
		{111, Many},
		{199, Many},
		{100, Other},
		{102, Other},
		{200, Other},
		{202, Other},
	}

	floatTests := []floatPluralTest{
		{0.1, Other},
		{0.2, Other},
		{0.3, Other},
		{1.1, Other},
		{1.2, Other},
		{1.3, Other},
		{2.1, Other},
		{2.2, Other},
		{2.3, Other},
		{3.1, Other},
		{3.2, Other},
		{3.3, Other},
	}

	language := LanguageWithCode("ar")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestEnglish(t *testing.T) {
	intTests := []intPluralTest{
		{0, Other},
		{1, One},
		{2, Other},
	}

	floatTests := []floatPluralTest{
		{0.1, Other},
		{0.2, Other},
		{1.1, Other},
		{1.2, Other},
		{2.1, Other},
		{2.2, Other},
	}

	language := LanguageWithCode("en")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestFrench(t *testing.T) {
	intTests := []intPluralTest{
		{0, One},
		{1, One},
		{2, Other},
	}

	floatTests := []floatPluralTest{
		{0.1, One},
		{0.2, One},
		{0.9, One},
		{1.1, One},
		{1.2, One},
		{1.9, One},
		{2.1, Other},
		{2.2, Other},
	}

	language := LanguageWithCode("fr")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestJapanese(t *testing.T) {
	intTests := []intPluralTest{
		{0, Other},
		{1, Other},
		{2, Other},
	}

	floatTests := []floatPluralTest{
		{0.1, Other},
		{0.2, Other},
		{1.1, Other},
		{1.2, Other},
		{2.1, Other},
		{2.2, Other},
	}

	language := LanguageWithCode("ja")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func TestSpanish(t *testing.T) {
	intTests := []intPluralTest{
		{0, Other},
		{1, One},
		{2, Other},
	}

	floatTests := []floatPluralTest{
		{0.1, Other},
		{0.2, Other},
		{1.1, Other},
		{1.2, Other},
		{2.1, Other},
		{2.2, Other},
	}

	language := LanguageWithCode("es")
	testInts(t, language, intTests)
	testIntsAsFloats(t, language, intTests)
	testFloats(t, language, floatTests)
}

func testInts(t *testing.T, language *Language, intTests []intPluralTest) {
	for _, test := range intTests {
		if pc := language.Int64PluralCategory(test.i); pc != test.pc {
			t.Errorf("Int64PluralCategory(%d) returned %s; expected %s", test.i, pc, test.pc)
		}
	}
}

func testIntsAsFloats(t *testing.T, language *Language, intTests []intPluralTest) {
	for _, test := range intTests {
		f := float64(test.i)
		if pc := language.Float64PluralCategory(f); pc != test.pc {
			t.Errorf("Float64PluralCategory(%f) returned %s; expected %s", f, pc, test.pc)
		}
	}
}

func testFloats(t *testing.T, language *Language, floatTests []floatPluralTest) {
	for _, test := range floatTests {
		if pc := language.Float64PluralCategory(test.f); pc != test.pc {
			t.Errorf("Float64PluralCategory(%f) returned %s; expected %s", test.f, pc, test.pc)
		}
	}
}
