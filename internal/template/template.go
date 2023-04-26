package template

import (
	"bytes"
	"text/template"

	"github.com/lindell/multi-gitter/internal/scm"
)

type templateStruct struct {
	FullName      string
	DefaultBranch string
}

func Template(pattern string, repository scm.Repository) (string, error) {
	t, err := template.New("stdout").Parse(pattern)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := t.Execute(&b, toTemplateStruct(repository)); err != nil {
		return "", err
	}

	return b.String(), nil
}

func toTemplateStruct(repo scm.Repository) *templateStruct {
	repo.FullName()
	repo.DefaultBranch()
	return &templateStruct{
		FullName:      repo.FullName(),
		DefaultBranch: repo.DefaultBranch(),
	}
}
