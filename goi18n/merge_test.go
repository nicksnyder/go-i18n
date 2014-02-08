package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestMergeExecute(t *testing.T) {
	resetDir(t, "testdata/output")
	files := []string{
		"testdata/input/en-US.one.json",
		"testdata/input/en-US.two.json",
		"testdata/input/fr-FR.json",
		"testdata/input/ar-AR.one.json",
		"testdata/input/ar-AR.two.json",
	}

	mc := &mergeCommand{
		translationFiles: files,
		sourceLocaleID:   "en-US",
		outdir:           "testdata/output",
		format:           "json",
	}
	if err := mc.execute(); err != nil {
		t.Fatal(err)
	}

	expectEqualFiles(t, "testdata/output/en-US.all.json", "testdata/expected/en-US.all.json")
	expectEqualFiles(t, "testdata/output/ar-AR.all.json", "testdata/expected/ar-AR.all.json")
	expectEqualFiles(t, "testdata/output/fr-FR.all.json", "testdata/expected/fr-FR.all.json")
	expectEqualFiles(t, "testdata/output/en-US.untranslated.json", "testdata/expected/en-US.untranslated.json")
	expectEqualFiles(t, "testdata/output/ar-AR.untranslated.json", "testdata/expected/ar-AR.untranslated.json")
	expectEqualFiles(t, "testdata/output/fr-FR.untranslated.json", "testdata/expected/fr-FR.untranslated.json")
}

func resetDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(dir, 0777); err != nil {
		t.Fatal(err)
	}
}

func expectEqualFiles(t *testing.T, expectedName, actualName string) {
	actual, err := ioutil.ReadFile(actualName)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile(expectedName)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Fatalf("contents of files did not match: %s, %s", expectedName, actualName)
	}
}
