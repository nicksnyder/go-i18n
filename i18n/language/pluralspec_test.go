package language

import (
	"fmt"
	"testing"
)

const onePlusEpsilon = "1.00000000000000000000000000000001"

func TestGetPluralSpec(t *testing.T) {
	tests := []struct {
		src  string
		spec *PluralSpec
	}{
		{"pl", pluralSpecs["pl"]},
		{"en", pluralSpecs["en"]},
		{"en-US", pluralSpecs["en"]},
		{"en_US", pluralSpecs["en"]},
		{"en-GB", pluralSpecs["en"]},
		{"zh-CN", pluralSpecs["zh"]},
		{"zh-TW", pluralSpecs["zh"]},
		{"pt-BR", pluralSpecs["pt-br"]},
		{"pt_BR", pluralSpecs["pt-br"]},
		{"pt-PT", pluralSpecs["pt"]},
		{"pt_PT", pluralSpecs["pt"]},
		{"zh-Hans-CN", pluralSpecs["zh"]},
		{"zh-Hant-TW", pluralSpecs["zh"]},
		{"zh-CN", pluralSpecs["zh"]},
		{"zh-TW", pluralSpecs["zh"]},
		{"zh-Hans", pluralSpecs["zh"]},
		{"zh-Hant", pluralSpecs["zh"]},
		{"en-US-en-US", pluralSpecs["en"]},
		{".en-US..en-US.", nil},
		{"zh, en-gb;q=0.8, en;q=0.7", nil},
		{"zh,en-gb;q=0.8,en;q=0.7", nil},
		{"xx, en-gb;q=0.8, en;q=0.7", nil},
		{"xx,en-gb;q=0.8,en;q=0.7", nil},
		{"xx-YY,xx;q=0.8,en-US,en;q=0.8,de;q=0.6,nl;q=0.4", nil},
		{"/foo/es/en.json", nil},
		{"xx-Yyen-US", nil},
		{"en US", nil},
		{"", nil},
		{"-", nil},
		{"_", nil},
		{".", nil},
		{"-en", nil},
		{"_en", nil},
		{"-en-", nil},
		{"_en_", nil},
		{"xx", nil},
	}
	for _, test := range tests {
		spec := getPluralSpec(test.src)
		if spec != test.spec {
			t.Errorf("getPluralSpec(%q) = %q expected %q", test.src, spec, test.spec)
		}
	}
}

type pluralTest struct {
	num    interface{}
	plural Plural
}

func TestArabic(t *testing.T ) {
	tests := []pluralTest{
		{0, Zero},
		{"0", Zero},
		{"0.0", Zero},
		{"0.00", Zero},
		{1, One},
		{"1", One},
		{"1.0", One},
		{"1.00", One},
		{onePlusEpsilon, Other},
		{2, Two},
		{"2", Two},
		{"2.0", Two},
		{"2.00", Two},
		{3, Few},
		{"3", Few},
		{"3.0", Few},
		{"3.00", Few},
		{10, Few},
		{"10", Few},
		{"10.0", Few},
		{"10.00", Few},
		{103, Few},
		{"103", Few},
		{"103.0", Few},
		{"103.00", Few},
		{110, Few},
		{"110", Few},
		{"110.0", Few},
		{"110.00", Few},
		{11, Many},
		{"11", Many},
		{"11.0", Many},
		{"11.00", Many},
		{99, Many},
		{"99", Many},
		{"99.0", Many},
		{"99.00", Many},
		{111, Many},
		{"111", Many},
		{"111.0", Many},
		{"111.00", Many},
		{199, Many},
		{"199", Many},
		{"199.0", Many},
		{"199.00", Many},
		{100, Other},
		{"100", Other},
		{"100.0", Other},
		{"100.00", Other},
		{102, Other},
		{"102", Other},
		{"102.0", Other},
		{"102.00", Other},
		{200, Other},
		{"200", Other},
		{"200.0", Other},
		{"200.00", Other},
		{202, Other},
		{"202", Other},
		{"202.0", Other},
		{"202.00", Other},
	}
	tests = appendFloatTests(tests, 0.1, 0.9, Other)
	tests = appendFloatTests(tests, 1.1, 1.9, Other)
	tests = appendFloatTests(tests, 2.1, 2.9, Other)
	tests = appendFloatTests(tests, 3.1, 3.9, Other)
	tests = appendFloatTests(tests, 4.1, 4.9, Other)
	runTests(t, "ar", tests)
}

func TestBelarusian(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{21, One},
		{11, Many},
		{52, Few},
		{101, One},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", One},
		{onePlusEpsilon, Other},
		{"2.0", Few},
		{"10.0", Many},
	}
	runTests(t, "be", tests)
}

func TestCatalan(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{"0", Other},
		{1, One},
		{"1", One},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{2, Other},
		{"2", Other},
	}
	tests = appendIntTests(tests, 2, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "ca", tests)
}

func TestChinese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "zh", tests)
}

func TestCzech(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{"0", Other},
		{1, One},
		{"1", One},
		{onePlusEpsilon, Many},
		{2, Few},
		{"2", Few},
		{3, Few},
		{"3", Few},
		{4, Few},
		{"4", Few},
		{5, Other},
		{"5", Other},
	}
	tests = appendFloatTests(tests, 0, 10, Many)
	runTests(t, "cs", tests)
}

func TestDanish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, One},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.1, 1.9, One)
	tests = appendFloatTests(tests, 2.0, 10.0, Other)
	runTests(t, "da", tests)
}

func TestDutch(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "nl", tests)
}

func TestEnglish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "en", tests)
}

func TestFrench(t *testing.T) {
	tests := []pluralTest{
		{0, One},
		{1, One},
		{onePlusEpsilon, One},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 1.9, One)
	tests = appendFloatTests(tests, 2.0, 10.0, Other)
	runTests(t, "fr", tests)
}

func TestGerman(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "de", tests)
}

func TestIcelandic(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{2, Other},
		{11, Other},
		{21, One},
		{111, Other},
		{"0.0", Other},
		{"0.1", One},
		{"2.0", Other},
	}
	runTests(t, "is", tests)
}

func TestItalian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "it", tests)
}

func TestJapanese(t *testing.T) {
	tests := appendIntTests(nil, 0, 10, Other)
	tests = appendFloatTests(tests, 0, 10, Other)
	runTests(t, "ja", tests)
}

func TestLithuanian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{2, Few},
		{3, Few},
		{9, Few},
		{10, Other},
		{11, Other},
		{"0.1", Many},
		{"0.7", Many},
		{"1.0", One},
		{onePlusEpsilon, Many},
		{"2.0", Few},
		{"10.0", Other},
	}
	runTests(t, "lt", tests)
}

func TestPolish(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{10, Many},
		{11, Many},
		{52, Few},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "pl", tests)
}

func TestPortuguese(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "pt", tests)
}

func TestPortugueseBrazilian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{"0.0", Other},
		{"0.1", One},
		{"0.01", One},
		{1, One},
		{"1", One},
		{"1.1", Other},
		{"1.01", Other},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 2.0, 10.0, Other)
	runTests(t, "pt-br", tests)
}

func TestRussian(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{21, One},
		{11, Many},
		{52, Few},
		{101, One},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "ru", tests)
}

func TestSpanish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{"1", One},
		{"1.0", One},
		{"1.00", One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 0.9, Other)
	tests = appendFloatTests(tests, 1.1, 10.0, Other)
	runTests(t, "es", tests)
}

func TestBulgarian(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{2, Other},
		{3, Other},
		{9, Other},
		{10, Other},
		{11, Other},
		{"0.1", Other},
		{"0.7", Other},
		{"1.0", One},
		{"1.001", Other},
		{onePlusEpsilon, Other},
		{"1.1", Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "bg", tests)
}

func TestSwedish(t *testing.T) {
	tests := []pluralTest{
		{0, Other},
		{1, One},
		{onePlusEpsilon, Other},
		{2, Other},
	}
	tests = appendFloatTests(tests, 0.0, 10.0, Other)
	runTests(t, "sv", tests)
}

func TestUkrainian(t *testing.T) {
	tests := []pluralTest{
		{0, Many},
		{1, One},
		{2, Few},
		{3, Few},
		{4, Few},
		{5, Many},
		{19, Many},
		{20, Many},
		{21, One},
		{11, Many},
		{52, Few},
		{101, One},
		{"0.1", Other},
		{"0.7", Other},
		{"1.5", Other},
		{"1.0", Other},
		{onePlusEpsilon, Other},
		{"2.0", Other},
		{"10.0", Other},
	}
	runTests(t, "uk", tests)
}

func appendIntTests(tests []pluralTest, from, to int, p Plural) []pluralTest {
	for i := from; i <= to; i++ {
		tests = append(tests, pluralTest{i, p})
	}
	return tests
}

func appendFloatTests(tests []pluralTest, from, to float64, p Plural) []pluralTest {
	stride := 0.1
	format := "%.1f"
	for f := from; f < to; f += stride {
		tests = append(tests, pluralTest{fmt.Sprintf(format, f), p})
	}
	tests = append(tests, pluralTest{fmt.Sprintf(format, to), p})
	return tests
}

func runTests(t *testing.T, specID string, tests []pluralTest) {
	spec := pluralSpecs[specID]
	for _, test := range tests {
		if plural, err := spec.Plural(test.num); plural != test.plural {
			t.Errorf("%s: PluralCategory(%#v) returned %s, %v; expected %s", specID, test.num, plural, err, test.plural)
		}
	}
}
