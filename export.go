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
		Title:       baseCfg.Site.Title,
		Link:        &feeds.Link{Href: fmt.Sprintf("https//%s.github.io/%s", baseCfg.Repository.Users[0], baseCfg.Repository.Name)},
		Description: baseCfg.Site.OneLineDesc,
		Author:      &feeds.Author{Name: baseCfg.Site.Author, Email: baseCfg.Site.Mail},
		Created:     now,
	}

	feed.Items = []*feeds.Item{}

	for _, issue := range issues {
		item := &feeds.Item{
			Title:       issue.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("https//%s.github.io/%s/%s", baseCfg.Repository.Users[0], baseCfg.Repository.Name, issue.Link)},
			Description: issue.Summary,
			Author:      &feeds.Author{Name: baseCfg.Site.Author, Email: baseCfg.Site.Mail},
			Created:     issue.Created,
		}
		feed.Items = append(feed.Items, item)
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(outDir, feedFile), []byte(atom), 0755)
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

func BuildSite(issues []*github.Issue, baseCfg BaseConfig, themeCfg ThemeConfig) error {
	exIssues, err := prepareIssues(issues)
	if err != nil {
		return err
	}

	for _, exis := range exIssues {
		err := exportIssue(exis, baseCfg, themeCfg)
		if err != nil {
			return err
		}
	}

	err = exportFeed(exIssues)
	if err != nil {
		return err
	}

	err = exportTemplate(themeCfg.IndexTemplate, exIssues, baseCfg, themeCfg)
	if err != nil {
		return err
	}

	return nil
}

func exportTemplate(template Template, issues []Issue, baseCfg BaseConfig, themeCfg ThemeConfig) error {
	data := map[string]interface{}{
		"Site":   baseCfg.Site,
		"Today":  time.Now(),
		"Issues": issues,
		"Custom": baseCfg.Custom,
		"Theme":  themeCfg,
	}

	var tmpl string
	var err error
	if template.Layout == "" {
		tmpl, err = mustache.RenderFile(filepath.Join(themeDir, template.Source), data)
	} else {
		tmpl, err = mustache.RenderFileInLayout(filepath.Join(themeDir, template.Source), filepath.Join(themeDir, template.Layout), data)
	}
	if err != nil {
		return err
	}

	outname := filepath.Join(outDir, template.Target)
	return ioutil.WriteFile(outname, []byte(tmpl), 0755)
}

func exportIssue(issue Issue, baseCfg BaseConfig, themeCfg ThemeConfig) error {
	data := map[string]interface{}{
		"Site":   baseCfg.Site,
		"Today":  time.Now(),
		"Issue":  issue,
		"Custom": baseCfg.Custom,
		"Theme":  themeCfg,
	}
	issueTmpl, err := mustache.RenderFileInLayout(filepath.Join(themeDir, "issue.mustache"), filepath.Join(themeDir, "layout.mustache"), data)
	if err != nil {
		return err
	}
	outname := filepath.Join(outDir, issue.Link)
	return ioutil.WriteFile(outname, []byte(issueTmpl), 0755)
}
