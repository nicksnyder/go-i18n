package language

import (
	"fmt"
	"testing"

	"github.com/nicksnyder/go-i18n/i18n/plural"
)

func TestParse(t *testing.T) {
	tests := []struct {
		src  string
		lang *Language
	}{
		{"en", languages["en"]},
		{"en-US", languages["en"]},
		{"en_US", languages["en"]},
		{"en-GB", languages["en"]},
		{"zh-CN", languages["zh"]},
		{"zh-TW", languages["zh"]},
		{"pt-BR", languages["pt-BR"]},
		{"pt_BR", languages["pt-BR"]},
		{"pt-PT", languages["pt"]},
		{"pt_PT", languages["pt"]},
		{"zh-Hans-CN", languages["zh"]},
		{"zh-Hant-TW", languages["zh"]},
		{"en-US-en-US", languages["en"]},
		{".en-US..en-US.", languages["en"]},
		{"zh, en-gb;q=0.8, en;q=0.7", languages["zh"]},
		{"zh,en-gb;q=0.8,en;q=0.7", languages["zh"]},
		{"xx, en-gb;q=0.8, en;q=0.7", languages["en"]},
		{"xx,en-gb;q=0.8,en;q=0.7", languages["en"]},
		{"xx-YY,xx;q=0.8,en-US,en;q=0.8,de;q=0.6,nl;q=0.4", languages["en"]},
		{"/foo/es/en.json", languages["en"]},
		{"xx-Yyen-US", nil},
		{"en US", nil},
		{"", nil},
		{"-", nil},
		{"_", nil},
		{"-en", nil},
		{"_en", nil},
		{"-en-", nil},
		{"_en_", nil},
		{"xx", nil},
	}
	for _, test := range tests {
		lang := Parse(test.src)
		if lang != test.lang {
			t.Errorf("Parse(%q) = %q expected %q", test.src, lang, test.lang)
		}
	}
}

type pluralTest struct {
	num interface{}
	pc  plural.Category
}

func TestArabic(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Zero},
		{"0", plural.Zero},
		{"0.0", plural.Zero},
		{"0.00", plural.Zero},
		{1, plural.One},
		{"1", plural.One},
		{"1.0", plural.One},
		{"1.00", plural.One},
		{2, plural.Two},
		{"2", plural.Two},
		{"2.0", plural.Two},
		{"2.00", plural.Two},
		{3, plural.Few},
		{"3", plural.Few},
		{"3.0", plural.Few},
		{"3.00", plural.Few},
		{10, plural.Few},
		{"10", plural.Few},
		{"10.0", plural.Few},
		{"10.00", plural.Few},
		{103, plural.Few},
		{"103", plural.Few},
		{"103.0", plural.Few},
		{"103.00", plural.Few},
		{110, plural.Few},
		{"110", plural.Few},
		{"110.0", plural.Few},
		{"110.00", plural.Few},
		{11, plural.Many},
		{"11", plural.Many},
		{"11.0", plural.Many},
		{"11.00", plural.Many},
		{99, plural.Many},
		{"99", plural.Many},
		{"99.0", plural.Many},
		{"99.00", plural.Many},
		{111, plural.Many},
		{"111", plural.Many},
		{"111.0", plural.Many},
		{"111.00", plural.Many},
		{199, plural.Many},
		{"199", plural.Many},
		{"199.0", plural.Many},
		{"199.00", plural.Many},
		{100, plural.Other},
		{"100", plural.Other},
		{"100.0", plural.Other},
		{"100.00", plural.Other},
		{102, plural.Other},
		{"102", plural.Other},
		{"102.0", plural.Other},
		{"102.00", plural.Other},
		{200, plural.Other},
		{"200", plural.Other},
		{"200.0", plural.Other},
		{"200.00", plural.Other},
		{202, plural.Other},
		{"202", plural.Other},
		{"202.0", plural.Other},
		{"202.00", plural.Other},
	}
	tests = appendFloatTests(tests, 0.1, 0.9, plural.Other)
	tests = appendFloatTests(tests, 1.1, 1.9, plural.Other)
	tests = appendFloatTests(tests, 2.1, 2.9, plural.Other)
	tests = appendFloatTests(tests, 3.1, 3.9, plural.Other)
	tests = appendFloatTests(tests, 4.1, 4.9, plural.Other)
	runTests(t, languages["ar"], tests)
}

func TestCatalan(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{"0", plural.Other},
		{1, plural.One},
		{"1", plural.One},
		{"1.0", plural.Other},
		{2, plural.Other},
		{"2", plural.Other},
	}
	tests = appendIntTests(tests, 2, 10, plural.Other)
	tests = appendFloatTests(tests, 0, 10, plural.Other)
	runTests(t, languages["ca"], tests)
}

func TestChinese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, plural.Other)
	tests = appendFloatTests(tests, 0, 10, plural.Other)
	runTests(t, languages["zh"], tests)
}

func TestCzech(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{"0", plural.Other},
		{1, plural.One},
		{"1", plural.One},
		{2, plural.Few},
		{"2", plural.Few},
		{3, plural.Few},
		{"3", plural.Few},
		{4, plural.Few},
		{"4", plural.Few},
		{5, plural.Other},
		{"5", plural.Other},
	}
	tests = appendFloatTests(tests, 0, 10, plural.Many)
	runTests(t, languages["cs"], tests)
}

func TestDanish(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.1, 1.9, plural.One)
	tests = appendFloatTests(tests, 2.0, 10.0, plural.Other)
	runTests(t, languages["da"], tests)
}

func TestDutch(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, plural.Other)
	runTests(t, languages["nl"], tests)
}

func TestEnglish(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, plural.Other)
	runTests(t, languages["en"], tests)
}

func TestFrench(t *testing.T) {
	tests := []pluralTest{
		{0, plural.One},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 1.9, plural.One)
	tests = appendFloatTests(tests, 2.0, 10.0, plural.Other)
	runTests(t, languages["fr"], tests)
}

func TestGerman(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, plural.Other)
	runTests(t, languages["de"], tests)
}

func TestItalian(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, plural.Other)
	runTests(t, languages["it"], tests)
}

func TestJapanese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, plural.Other)
	tests = appendFloatTests(tests, 0, 10, plural.Other)
	runTests(t, languages["ja"], tests)
}

func TestLithuanian(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Few},
		{3, plural.Few},
		{9, plural.Few},
		{10, plural.Other},
		{11, plural.Other},
		{"0.1", plural.Many},
		{"0.7", plural.Many},
		{"1.0", plural.One},
		{"2.0", plural.Few},
		{"10.0", plural.Other},
	}
	runTests(t, languages["lt"], tests)
}

func TestPortuguese(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, plural.Other)
	runTests(t, languages["pt"], tests)
}

func TestPortugueseBrazilian(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{"0.0", plural.Other},
		{"0.1", plural.One},
		{"0.01", plural.One},
		{1, plural.One},
		{"1", plural.One},
		{"1.1", plural.Other},
		{"1.01", plural.Other},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 2.0, 10.0, plural.Other)
	runTests(t, languages["pt-BR"], tests)
}

func TestSpanish(t *testing.T) {
	tests := []pluralTest{
		{0, plural.Other},
		{1, plural.One},
		{"1", plural.One},
		{"1.0", plural.One},
		{"1.00", plural.One},
		{2, plural.Other},
	}
	tests = appendFloatTests(tests, 0.0, 0.9, plural.Other)
	tests = appendFloatTests(tests, 1.1, 10.0, plural.Other)
	runTests(t, languages["es"], tests)
}

func appendIntTests(tests []pluralTest, from, to int, pc plural.Category) []pluralTest {
	for i := from; i <= to; i++ {
		tests = append(tests, pluralTest{i, pc})
	}
	return tests
}

func appendFloatTests(tests []pluralTest, from, to float64, pc plural.Category) []pluralTest {
	stride := 0.1
	format := "%.1f"
	for f := from; f < to; f += stride {
		tests = append(tests, pluralTest{fmt.Sprintf(format, f), pc})
	}
	tests = append(tests, pluralTest{fmt.Sprintf(format, to), pc})
	return tests
}

func runTests(t *testing.T, language *Language, tests []pluralTest) {
	for _, test := range tests {
		if pc, err := language.PluralCategory(test.num); pc != test.pc {
			t.Errorf("%s: PluralCategory(%#v) returned %s, %v; expected %s", language.ID, test.num, pc, err, test.pc)
		}
	}
}
