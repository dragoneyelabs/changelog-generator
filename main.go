package main

import (
	"github.com/dragoneyelabs/changelog-generator/changelog"

	"flag"
	"path"
)

var pathFlag = flag.String("path", "", "Git project path")

func main() {
	flag.Parse()
	if pathFlag == nil || *pathFlag == "" {
		panic("a path for the git project must be specified")
	}

	if err := changelog.Generate(*pathFlag, path.Join(*pathFlag, "CHANGELOG.md")); err != nil {
		panic(err)
	}
}
