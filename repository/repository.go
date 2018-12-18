package repository

import (
	"errors"
	"strings"
)

const (
	// Bitbucket parser type
	Bitbucket Type = "bitbucket"
	// Github parser type
	Github Type = "github"
	// Unknown parser type
	Unknown Type = "unknown"
)

// Type represents the parser type
type Type string

// URLParser allows to return repository URLs
type URLParser interface {
	GetCommitURL(hash string) string
	GetTagURL(tag string) string
}

// GetType returns the repository type from a string. Returns Unknown type if not found
func GetType(repositoryType string) Type {
	switch Type(strings.ToLower(repositoryType)) {
	case Bitbucket:
		return Bitbucket
	case Github:
		return Github
	default:
		return Unknown
	}
}

// GetParser returns the URLParser for the given repository type and url
func GetParser(repositoryType Type, url string) (URLParser, error) {
	switch repositoryType {
	case Bitbucket:
		return newBitbucketParser(url)
	case Github:
		return newGithubParser(url)
	default:
		return nil, errors.New("invalid parser")
	}
}
