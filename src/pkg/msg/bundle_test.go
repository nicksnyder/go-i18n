package msg

import (
	"sort"
	"testing"
)

func TestMessages(t *testing.T) {
	m := []Message{
		Message{"", "", "translated", ""},
		Message{"", "", "translated", "t1"},
		Message{"", "", "translated", "t2"},
		Message{"", "", "translated", ""},
		Message{"", "", "untranslated", ""},
		Message{"", "", "untranslated", ""},
	}
	expectedAll := []Message{
		Message{"0322ba425026b83b6a285b243affe34e", "", "translated", "t2"},
		Message{"f0e4c362bb622adba484f6e29f0eef4f", "", "untranslated", ""},
	}
	expectedTranslated := []Message{
		Message{"0322ba425026b83b6a285b243affe34e", "", "translated", "t2"},
	}
	expectedUntranslated := []Message{
		Message{"f0e4c362bb622adba484f6e29f0eef4f", "", "untranslated", ""},
	}

	b := NewBundle()
	b.AddMessages(m)
	mustEqual(t, expectedAll, b.Messages())
	mustEqual(t, expectedTranslated, b.TranslatedMessages())
	mustEqual(t, expectedUntranslated, b.UntranslatedMessages())
}

func TestSortByContent(t *testing.T) {
	m := []Message{
		Message{"", "y", "b", "y"},
		Message{"", "w", "d", "w"},
		Message{"", "x", "c", "x"},
		Message{"", "z", "a", "z"},
	}
	expected := []Message{
		Message{"959848ca10cc8a60da818ac11523dc63", "z", "a", "z"},
		Message{"0b8dd1a01a737950a643a578f08f7900", "y", "b", "y"},
		Message{"831c4baa8a44083a6434b892d573846b", "x", "c", "x"},
		Message{"dd3ba2cca7da8526615be73d390527ac", "w", "d", "w"},
	}

	b := NewBundle()
	b.AddMessages(m)
	sort.Sort(b)
	mustEqual(t, expected, b.Messages())
}

func TestSortByContext(t *testing.T) {
	m := []Message{
		Message{"", "y", "a", "m"},
		Message{"", "w", "a", "o"},
		Message{"", "x", "a", "n"},
		Message{"", "z", "a", "l"},
	}
	expected := []Message{
		Message{"c68c559d956d4ca20f435ed74a6e71e6", "w", "a", "o"},
		Message{"53e59fface936ea788f7cf51e7b25531", "x", "a", "n"},
		Message{"d74600e380dbf727f67113fd71669d88", "y", "a", "m"},
		Message{"959848ca10cc8a60da818ac11523dc63", "z", "a", "l"},
	}

	b := NewBundle()
	b.AddMessages(m)
	sort.Sort(b)
	mustEqual(t, expected, b.Messages())
}

func mustEqual(t *testing.T, expected, actual []Message) {
	if ok := equal(expected, actual); !ok {
		t.Fatalf("Slices not equal!\nExpected:\n%#v\nActual:\n%#v", expected, actual)
	}
}

func equal(a, b []Message) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Id != b[i].Id ||
			a[i].Context != b[i].Context ||
			a[i].Content != b[i].Content ||
			a[i].Translation != b[i].Translation {
			return false
		}
	}
	return true
}
