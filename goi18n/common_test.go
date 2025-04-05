package main

import (
	"os"
	"testing"
)

func mustTempDir(prefix string) string {
	outdir, err := os.MkdirTemp("", prefix)
	if err != nil {
		panic(err)
	}
	return outdir
}

func mustRemoveAll(t *testing.T, path string) {
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
}
