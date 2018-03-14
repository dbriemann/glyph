package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/gosimple/slug"
	"github.com/microcosm-cc/bluemonday"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type Issue struct {
	Title   string
	Link    string // slugified title
	Content template.HTML
	Summary template.HTML
}

func prepareIssues(issues []*github.Issue) ([]Issue, error) {
	export := []Issue{}

	for _, issue := range issues {
		exIssue := Issue{
			Title: issue.GetTitle(),
			Link:  slug.Make(issue.GetTitle()) + ".html",
		}
		if exIssue.Title != "" {
			unsafe := blackfriday.Run([]byte(issue.GetBody()))
			html := string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
			exIssue.Content = template.HTML(html)
			if len(exIssue.Content) > 200 {
				exIssue.Summary = exIssue.Content[:200]
			} else {
				exIssue.Summary = exIssue.Content
			}
			export = append(export, exIssue)
		}
		// We ignore issues with empty titles.
	}

	return export, nil
}

func BuildSite(issues []*github.Issue, cfg *Config) error {
	exIssues, err := prepareIssues(issues)
	if err != nil {
		return err
	}

	for _, exis := range exIssues {
		err := exportIssue(exis, cfg)
		if err != nil {
			return err
		}
	}

	err = exportIndex(exIssues, cfg)
	if err != nil {
		return err
	}

	return nil
}

func exportIssue(issue Issue, cfg *Config) error {
	raw, err := ioutil.ReadFile("tmpl/issue.tmpl")
	if err != nil {
		return err
	}
	itmpl, err := template.New(issue.Title).Parse(string(raw))
	if err != nil {
		return err
	}

	data := struct {
		TITLE         string
		ISSUE_TITLE   string
		ISSUE_CONTENT template.HTML
	}{
		TITLE:         cfg.Site.Title,
		ISSUE_TITLE:   issue.Title,
		ISSUE_CONTENT: issue.Content,
	}

	f, err := os.Create(filepath.Join(cfg.Repository.TargetDir, issue.Link))
	if err != nil {
		return err
	}
	defer f.Close()

	err = itmpl.Execute(f, data)
	return err
}

func exportIndex(issues []Issue, cfg *Config) error {
	raw, err := ioutil.ReadFile("tmpl/index.tmpl")
	if err != nil {
		return err
	}
	index, err := template.New("index").Parse(string(raw))
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
