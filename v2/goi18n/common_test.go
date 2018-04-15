package main

import "io/ioutil"

func mustTempDir(prefix string) string {
	outdir, err := ioutil.TempDir("", "TestExtractCommand")
	if err != nil {
		panic(err)
	}
	return outdir
}
