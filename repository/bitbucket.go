package repository

import (
	"github.com/dragoneyelabs/changelog-generator/utils"

	"fmt"
	"strings"
)

type bitbucketParser struct {
	Host       string
	Project    string
	Repository string
}

func (p *bitbucketParser) GetCommitURL(hash string) string {
	return fmt.Sprintf("https://%s/bitbucket/projects/%s/repos/%s/commits/%s", p.Host, p.Project, p.Repository, hash)
}

func (p *bitbucketParser) GetTagURL(tag string) string {
	return fmt.Sprintf("https://%s/bitbucket/projects/%s/repos/%s/commits?until=%s",
		p.Host, p.Project, p.Repository, tag)
}

func newBitbucketParser(url string) (*bitbucketParser, error) {
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

	return &bitbucketParser{
		Host:       host,
		Project:    s[0],
		Repository: s[1],
	}, nil
}
