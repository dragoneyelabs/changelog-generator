package template

import (
	"bytes"
	"os"
	"text/template"
	"time"
)

const (
	mainTmplStr    = "{{range .Tags}}{{if or .Feats .Changes .Fixes .Tests .Docs}}{{template \"tag\" .}}\n\n{{end}}{{end}}"
	commitsTmplStr = "{{range .}}{{if .Msg}}- {{.Msg}} (" +
		"{{if .URL}}[{{printf \"%.7s\" .Hash}}]({{.URL}})" +
		"{{else}}{{printf \"%.7s\" .Hash}}" +
		"{{end}}){{end}}\n{{end}}"
	//nolint:lll
	tagTmplStr = `# {{if .Name}}{{if .URL}}[{{.Name}}]({{.URL}}){{else}}{{.Name}}{{end}}{{else}}Unreleased{{end}} ({{.Date.Format "2006-01-02"}}){{if .Feats}}

### Features

{{template "commits" .Feats}}{{end}}{{if .Changes}}

### Changes

{{template "commits" .Changes}}{{end}}{{if .Fixes}}

### Bug Fixes

{{template "commits" .Fixes}}{{end}}{{if .Tests}}

### Tests

{{template "commits" .Tests}}{{end}}{{if .Docs}}

### Docs

{{template "commits" .Docs}}{{end}}`
)

var t *template.Template

// Params stores the template params
type Params struct {
	Tags []TagParam
}

// TagParam stores the tag template params
type TagParam struct {
	Name    string
	URL     string
	Date    time.Time
	Feats   []CommitParam
	Changes []CommitParam
	Fixes   []CommitParam
	Tests   []CommitParam
	Docs    []CommitParam
}

// CommitParam stores the commit template params
type CommitParam struct {
	Msg  string
	Hash string
	URL  string
}

func init() {
	t = template.Must(template.New("main").Parse(mainTmplStr))
	template.Must(t.New("tag").Parse(tagTmplStr))
	template.Must(t.New("commits").Parse(commitsTmplStr))
}

// WriteFile writes the template in the filePath using the given template params
func WriteFile(filePath string, params Params) error {
	var b bytes.Buffer
	err := t.Execute(&b, params)
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	_, err = f.Write(b.Bytes())
	return err
}
