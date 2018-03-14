package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/gosimple/slug"
	"github.com/hoisie/mustache"
	gfm "github.com/shurcooL/github_flavored_markdown"
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
			Title: issue.GetTitle(),
			Link:  slug.Make(issue.GetTitle()) + ".html",
		}
		if exIssue.Title != "" {
			exIssue.Content = string(gfm.Markdown([]byte(issue.GetBody())))
			if len(exIssue.Content) > 200 { // TODO
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
	data := struct {
		TITLE         string
		ISSUE_TITLE   string
		ISSUE_CONTENT string
	}{
		TITLE:         cfg.Site.Title,
		ISSUE_TITLE:   issue.Title,
		ISSUE_CONTENT: issue.Content,
	}

	issueTmpl := mustache.RenderFileInLayout("tmpl/issue.mustache", "tmpl/layout.mustache", data)
	outname := filepath.Join(cfg.Repository.TargetDir, issue.Link)
	err := ioutil.WriteFile(outname, []byte(issueTmpl), 0755)
	return err
}

func exportIndex(issues []Issue, cfg *Config) error {
	data := struct {
		TITLE  string
		ISSUES []Issue
	}{
		TITLE:  cfg.Site.Title,
		ISSUES: issues,
	}
	indexTmpl := mustache.RenderFileInLayout("tmpl/index.mustache", "tmpl/layout.mustache", data)
	outname := filepath.Join(cfg.Repository.TargetDir, "index.html")
	err := ioutil.WriteFile(outname, []byte(indexTmpl), 0755)
	return err
}
