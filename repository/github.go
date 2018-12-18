package repository

import (
	"github.com/dragoneyelabs/changelog-generator/utils"

	"fmt"
	"strings"
)

type githubParser struct {
	Host       string
	Owner      string
	Repository string
}

func (p *githubParser) GetCommitURL(hash string) string {
	return fmt.Sprintf("https://%s/%s/%s/commit/%s", p.Host, p.Owner, p.Repository, hash)
}

// https://github.com/git-chglog/git-chglog.git
func (p *githubParser) GetTagURL(tag string) string {
	return fmt.Sprintf("https://%s/%s/%s/commits/%s", p.Host, p.Owner, p.Repository, tag)
}

func newGithubParser(url string) (*githubParser, error) {
	host, err := utils.GetHost(url)
	if err != nil {
		return nil, err
	}
	path, err := utils.GetPath(url)
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, ".git")
	path = strings.TrimPrefix(path, "/")
	s := strings.SplitN(path, "/", 2)

	return &githubParser{
		Host:       host,
		Owner:      s[0],
		Repository: s[1],
	}, nil
}
