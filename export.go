package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cbroglie/mustache"
	"github.com/google/go-github/github"
	"github.com/gorilla/feeds"
	"github.com/gosimple/slug"

	gfm "github.com/shurcooL/github_flavored_markdown"
)

type Issue struct {
	Title   string
	Link    string // slugified title
	Content string
	Summary string
	Labels  []string
	Created time.Time
}

func exportFeed(issues []Issue) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       cfg.Site.Title,
		Link:        &feeds.Link{Href: fmt.Sprintf("https//%s.github.io/%s", cfg.Repository.Users[0], cfg.Repository.Name)},
		Description: cfg.Site.OneLineDesc,
		Author:      &feeds.Author{Name: cfg.Site.Author, Email: cfg.Site.Mail},
		Created:     now,
	}

	feed.Items = []*feeds.Item{}

	for _, issue := range issues {
		item := &feeds.Item{
			Title:       issue.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https//%s.github.io/%s/%s", cfg.Repository.Users[0], cfg.Repository.Name, issue.Link)},
			Description: issue.Summary,
			Author:      &feeds.Author{Name: cfg.Site.Author, Email: cfg.Site.Mail},
			Created:     issue.Created,
		}
		feed.Items = append(feed.Items, item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(cfg.Repository.OutputDir, feedFile), []byte(atom), 0755)
}

func prepareIssues(issues []*github.Issue) ([]Issue, error) {
	export := []Issue{}

	for _, issue := range issues {
		exIssue := Issue{
			Title:   issue.GetTitle(),
			Link:    slug.Make(issue.GetTitle()) + ".html",
			Created: issue.GetCreatedAt(),
			Labels:  []string{},
		}
		// TODO maybe add syntax highlighting with chroma here?
		if exIssue.Title != "" {
			exIssue.Content = string(gfm.Markdown([]byte(issue.GetBody())))
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(exIssue.Content))
			if err == nil {
				// Use first paragraph(p) as summary.
				firstp := doc.Find("p").First()
				html, err := firstp.Html()
				if err == nil {
					exIssue.Summary = html
				}
			}

			for _, label := range issue.Labels {
				exIssue.Labels = append(exIssue.Labels, label.GetName())
			}

			export = append(export, exIssue)
		}
		// We ignore issues with empty titles.
	}

	return export, nil
}

func BuildSite(issues []*github.Issue, cfg Config) error {
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

	err = exportFeed(exIssues)
	if err != nil {
		return err
	}

	err = exportAbout(cfg)
	if err != nil {
		return err
	}

	err = exportIndex(exIssues, cfg)
	if err != nil {
		return err
	}

	return nil
}

func exportAbout(cfg Config) error {
	data := map[string]interface{}{
		"Site":   cfg.Site,
		"Today":  time.Now(),
		"Custom": cfg.Custom,
	}
	indexTmpl, err := mustache.RenderFileInLayout(filepath.Join(themeDir, "about.mustache"), filepath.Join(themeDir, "layout.mustache"), data)
	if err != nil {
		return err
	}
	outname := filepath.Join(cfg.Repository.OutputDir, "about.html")
	return ioutil.WriteFile(outname, []byte(indexTmpl), 0755)
}

func exportIssue(issue Issue, cfg Config) error {
	data := map[string]interface{}{
		"Site":   cfg.Site,
		"Today":  time.Now(),
		"Issue":  issue,
		"Custom": cfg.Custom,
	}
	issueTmpl, err := mustache.RenderFileInLayout(filepath.Join(themeDir, "issue.mustache"), filepath.Join(themeDir, "layout.mustache"), data)
	if err != nil {
		return err
	}
	outname := filepath.Join(cfg.Repository.OutputDir, issue.Link)
	return ioutil.WriteFile(outname, []byte(issueTmpl), 0755)
}

func exportIndex(issues []Issue, cfg Config) error {
	data := map[string]interface{}{
		"Site":   cfg.Site,
		"Today":  time.Now(),
		"Issues": issues,
		"Custom": cfg.Custom,
	}
	indexTmpl, err := mustache.RenderFileInLayout(filepath.Join(themeDir, "index.mustache"), filepath.Join(themeDir, "layout.mustache"), data)
	if err != nil {
		return err
	}
	outname := filepath.Join(cfg.Repository.OutputDir, "index.html")
	return ioutil.WriteFile(outname, []byte(indexTmpl), 0755)
}
