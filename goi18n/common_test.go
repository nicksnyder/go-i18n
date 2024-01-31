package main

import "os"

func mustTempDir(prefix string) string {
	outdir, err := os.MkdirTemp("", prefix)
	if err != nil {
		panic(err)
	}
	return outdir
}
