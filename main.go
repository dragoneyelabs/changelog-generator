package main

import (
	"github.com/dragoneyelabs/changelog-generator/changelog"

	"flag"
	"os"
	"path"
)

var pathFlag = flag.String("path", "", "Git project path")

func main() {
	flag.Parse()
	if pathFlag == nil || *pathFlag == "" {
		panic("a path for the git project must be specified")
	}

	w, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := changelog.Generate(*pathFlag, path.Join(w, "CHANGELOG.md")); err != nil {
		panic(err)
	}
}
