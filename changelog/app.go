package changelog

import (
	"github.com/dragoneyelabs/changelog-generator/git"
	"github.com/dragoneyelabs/changelog-generator/repository"
	"github.com/dragoneyelabs/changelog-generator/template"
	"github.com/dragoneyelabs/changelog-generator/utils"

	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	defaultKey = "DEFAULT"
	featsKey   = "feat:"
	changesKey = "refactor:"
	fixesKey   = "fix:"
	testsKey   = "test:"
	docsKey    = "docs:"
)

var prefixes = []string{
	featsKey,
	changesKey,
	fixesKey,
	testsKey,
	docsKey,
}

// Generate executes the changelog generation
func Generate(gitPath, changelogFilePath, repositoryType string) error {
	r, err := git.GetLocalRepository(gitPath)
	if err != nil {
		return err
	}
	head, err := r.GetHead()
	if err != nil {
		return err
	}
	tags, err := r.GetTagsSortedByDateAsc() // TODO: get remote tags
	if err != nil {
		return err
	}

	var parser repository.URLParser
	var urls []string
	if urls, err = r.GetRemoteURLs("origin"); err == nil && len(urls) > 0 { // TODO: pass remote as parameter
		parser, err = repository.GetParser(repository.GetType(repositoryType), urls[0])
		if err != nil {
			fmt.Printf("WARNING: %s\n", err)
		}
	}

	params, err := buildTemplateParams(r, head, tags, parser)
	if err != nil {
		return err
	}

	return template.WriteFile(changelogFilePath, *params)
}

func buildTemplateParams(r *git.Repository, head *git.Commit, tags []git.Tag,
	parser repository.URLParser) (*template.Params, error) {
	params := template.Params{}
	seen := make([]git.Hash, 0)
	for i := range tags {
		commits, err := r.GetCommitsSortedByDateDesc(tags[i].Hash, seen...)
		if err != nil {
			return nil, err
		}

		seen = append(seen, getHashes(commits)...)
		params.Tags = append(params.Tags, buildTagParam(tags[i].Name, tags[i].Date, commits, parser))
	}

	if len(tags) > 0 && head.Hash != tags[0].Hash || len(tags) == 0 {
		commits, err := r.GetCommitsSortedByDateDesc(head.Hash, seen...)
		if err != nil {
			return nil, err
		}
		params.Tags = append(params.Tags, buildTagParam("", head.Date, commits, parser))
	}

	// order by date desc
	sort.Slice(params.Tags, func(i, j int) bool {
		return params.Tags[j].Date.Before(params.Tags[i].Date)
	})

	return &params, nil
}

func getHashes(commits []git.Commit) []git.Hash {
	h := make([]git.Hash, 0)
	for i := range commits {
		h = append(h, commits[i].Hash)
	}
	return h
}

func buildTagParam(name string, date time.Time, commits []git.Commit, parser repository.URLParser) template.TagParam {
	commitsMap := mapCommitsByPrefix(commits)
	url := ""
	if parser != nil {
		url = parser.GetTagURL(name)
	}
	return template.TagParam{
		Name:    name,
		Date:    date,
		URL:     url,
		Feats:   toCommitParams(commitsMap[featsKey], featsKey, parser),
		Changes: toCommitParams(commitsMap[changesKey], changesKey, parser),
		Fixes:   toCommitParams(commitsMap[fixesKey], fixesKey, parser),
		Tests:   toCommitParams(commitsMap[testsKey], testsKey, parser),
		Docs:    toCommitParams(commitsMap[docsKey], docsKey, parser),
	}
}

func mapCommitsByPrefix(commits []git.Commit) map[string][]git.Commit {
	m := make(map[string][]git.Commit)
	for i := range commits {
		prefix := utils.GetPrefix(strings.ToLower(strings.TrimSpace(commits[i].Message)), prefixes, defaultKey)
		m[prefix] = append(m[prefix], commits[i])
	}
	return m
}

func toCommitParams(commits []git.Commit, prefix string, parser repository.URLParser) []template.CommitParam {
	params := make([]template.CommitParam, 0)
	msgArr := make([]string, 0)
	for i := range commits {
		msg := strings.TrimSpace(strings.TrimPrefix(commits[i].Message, prefix))
		if msg == "" {
			continue
		}
		msg = utils.Capitalize(strings.TrimSpace(strings.Split(msg, "\n")[0]))

		// avoid duplicates
		if utils.InArray(msg, msgArr) {
			continue
		}
		msgArr = append(msgArr, msg)

		url := ""
		if parser != nil {
			url = parser.GetCommitURL(commits[i].Hash.String())
		}

		params = append(params, template.CommitParam{
			Msg:  msg,
			Hash: commits[i].Hash.String(),
			URL:  url,
		})
	}
	return params
}
