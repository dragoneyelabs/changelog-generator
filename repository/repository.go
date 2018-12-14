package repository

import (
	"errors"
)

const (
	// Bitbucket parser type
	Bitbucket Type = iota
)

// Type represents the parser type
type Type int

// URLParser allows to return repository URLs
type URLParser interface {
	GetCommitURL(hash string) string
	GetTagURL(tag string) string
}

// GetParser returns the URLParser for the given repository type and url
func GetParser(repositoryType Type, url string) (URLParser, error) {
	switch repositoryType {
	case Bitbucket:
		return newBitbucketParser(url)
	default:
		return nil, errors.New("invalid parser")
	}
}
