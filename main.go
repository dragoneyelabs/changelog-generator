package main

import (
	"github.com/dragoneyelabs/changelog-generator/changelog"

	"flag"
	"path"
)

var (
	pathFlag           = flag.String("path", "", "Git project path")
	repositoryTypeFlag = flag.String("repo", "", "Repository Type (bitbucket, github)")
)

func main() {
	flag.Parse()
	if pathFlag == nil || *pathFlag == "" {
		panic("a path for the git project must be specified")
	}
	repositoryType := ""
	if repositoryTypeFlag != nil {
		repositoryType = *repositoryTypeFlag
	}

	if err := changelog.Generate(*pathFlag, path.Join(*pathFlag, "CHANGELOG.md"), repositoryType); err != nil {
		panic(err)
	}
}
