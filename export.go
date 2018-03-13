package main

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/dbriemann/juicer/tmpl"
	"github.com/google/go-github/github"
	"github.com/gosimple/slug"
)

type Issue struct {
	Title   string
	Link    string // slugified title
	Content string
	Summary string
}

func prepareIssues(issues []*github.Issue) ([]Issue, error) {
	export := []Issue{}

	for _, issue := range issues {
		exIssue := Issue{
			Title:   issue.GetTitle(),
			Link:    slug.Make(issue.GetTitle()),
			Content: issue.GetBody(),
		}
		if exIssue.Title != "" {
			// We ignore issues with empty titles.
			export = append(export, exIssue)
		}
	}

	return export, nil
}

func BuildSite(issues []*github.Issue, cfg *Config) error {
	exIssues, err := prepareIssues(issues)
	if err != nil {
		return err
	}

	err = exportIndex(exIssues, cfg)
	if err != nil {
		return err
	}

	return nil
}

func exportIndex(issues []Issue, cfg *Config) error {
	// TODO copy base files

	index, err := template.New("name").Parse(tmpl.IndexTmpl)
	if err != nil {
		return err
	}

	data := struct {
		TITLE  string
		ISSUES []Issue
	}{
		TITLE:  cfg.Site.Title,
		ISSUES: issues,
	}

	f, err := os.Create(filepath.Join(cfg.Repository.TargetDir, "index.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	err = index.Execute(f, data)
	return err
}
