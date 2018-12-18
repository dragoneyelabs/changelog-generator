# Changelog Generator

A simple changelog generator for Git projects, implemented in Go.

## Getting Started

### Installation

#### Go

```bash
$ go get -u github.com/dragoneyelabs/changelog-generator
```

#### Source

Clone the Github repository into your local machine and build the binary:

```bash
$ make build [linux=1] [darwin=1] [windows=1]
```

or

```bash
$ GOARCH=<386, amd64, arm, ...> GOOS=<linux, darwin, windows, ...> go build -o bin/changelog-generator github.com/dragoneyelabs/changelog-generator
```

### Usage

In order to use the `changelog-generator`, it is needed the path to a Git project.

Optionally, we can indicate the repository type as well. Only `bitbucket` and `github` types are currently supported.

#### Go

```bash
$ go run *.go -path=<project path> -repo=<bitbucket,github>
```

#### Binary

```bash
$ changelog-generator -path=<project path> -repo=<bitbucket,github>
```

#### Docker

Alternatively, we can run the changelog generator from a Docker container:

```bash
$ docker run -v `pwd`:/changelog pruizpar/changelog-generator:master -path=/changelog -repo=github
```
