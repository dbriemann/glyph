package main

import (
	"context"
	"fmt"
	"path/filepath"

	"os"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	feedFile   = "feed.atom"
	baseDir    = "themes"
	configFile = "config.toml"

	outDir = "docs"
)

var (
	baseCfg  BaseConfig
	themeCfg ThemeConfig
	themeDir = ""
)

// TODO add syntax highlighting with chroma?
// TODO find repo internal issue links in issues and replace them

func main() {
	// Read config file and check sanity.
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		bye(fmt.Sprintf("could not read config file: %s", err.Error()), 1)
	}
	if _, err := toml.Decode(string(raw), &baseCfg); err != nil {
		bye(fmt.Sprintf("could not parse config file: %s", err.Error()), 1)
	}

	// Fetch Github access token from environment.
	token := os.Getenv(baseCfg.GithubToken)

	// Test config data sanity.
	if baseCfg.Repository.Name == "" {
		bye("config file: repository name missing", 1)
	}
	if len(baseCfg.Repository.Users) < 1 {
		bye("config file: no user(s) provided", 1)
	}
	if baseCfg.Site.Title == "" {
		fmt.Println("warning: no site title set")
	}
	if baseCfg.Site.Author == "" {
		fmt.Println("warning: no site author set")
	}
	if baseCfg.Site.OneLineDesc == "" {
		fmt.Println("warning: no one line description set")
	}
	if baseCfg.Site.Mail == "" {
		fmt.Println("warning: no contact mail set")
	}
	if baseCfg.Site.Twitter == "" {
		fmt.Println("warning: no twitter handle set")
	}
	if baseCfg.Site.Theme == "" {
		fmt.Println("warning: no theme set.. falling back to default")
		baseCfg.Site.Theme = "default"
	}

	// Set output directory.
	finfo, err := os.Stat(outDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(outDir, 0755)
			if err != nil {
				bye(fmt.Sprintf("could not create output directory: %s", err.Error()), 1)
			}

			err = ioutil.WriteFile(filepath.Join(outDir, ".nojekyll"), []byte(""), 0755)
			if err != nil {
				bye(fmt.Sprintf("could not write .nojekyll file: %s", err.Error()), 1)
			}
		} else {
			bye(fmt.Sprintf("could not access output directory: %s", err.Error()), 1)
		}
	} else {
		if !finfo.IsDir() {
			bye(fmt.Sprintf("%s should be a directory but is a file", outDir), 1)
		}
	}

	// Try to open theme folder..
	themeDir = filepath.Join(baseDir, baseCfg.Site.Theme)
	finfo, err = os.Stat(themeDir)
	if err != nil || !finfo.IsDir() {
		fmt.Sprintf("cannot load theme %s: %s\n", baseCfg.Site.Theme, err.Error())
		fmt.Println("falling back to default theme")
		themeDir = filepath.Join(baseDir, "default")
	}

	// Read theme config file
	raw, err = ioutil.ReadFile(filepath.Join(themeDir, configFile))
	if err != nil {
		bye(fmt.Sprintf("could not read theme config file: %s", err.Error()), 1)
	}
	if _, err := toml.Decode(string(raw), &themeCfg); err != nil {
		bye(fmt.Sprintf("could not parse config file: %s", err.Error()), 1)
	}

	// Test config data sanity.
	if themeCfg.Name == "" {
		bye("theme config file: no name provided", 1)
	}
	if themeCfg.IndexTemplate.Source == "" {
		bye("theme config file: no index template source provided", 1)
	}
	if themeCfg.IndexTemplate.Target == "" {
		bye("theme config file: no index template target provided", 1)
	}
	if themeCfg.IssueTemplate.Source == "" {
		bye("theme config file: no issue template source provided", 1)
	}

	for _, t := range themeCfg.OtherTemplates {
		if t.Source == "" || t.Target == "" {
			bye("theme config file: custom template data incomplete (source or target missing)", 1)
		}
	}

	// Copy theme files except mustache template files and theme config file.
	files, err := ioutil.ReadDir(themeDir)
	if err != nil {
		bye(fmt.Sprintf("could not read access directory: %s", err.Error()), 1)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".mustache" && f.Name() != configFile && f.Name() != "README.md" {
			src := filepath.Join(themeDir, f.Name())
			dst := filepath.Join(outDir, f.Name())
			if err := copyFile(src, dst); err != nil {
				bye(fmt.Sprintf("could not copy file: %s", err.Error()), 1)
			}
		}
	}

	ctx := context.Background()
	var client *github.Client

	if token == "" {
		// No Github token. Create client without authing to API.
		client = github.NewClient(nil)
	} else {
		// Found a token. Auth and create client.
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	if client == nil {
		bye("client dead", 1)
	}

	issues, err := FetchIssues(client, ctx, baseCfg)
	if err != nil {
		bye(fmt.Sprintf("could not fetch issues from github: %s", err.Error()), 1)
	}

	err = BuildSite(issues, baseCfg, themeCfg)
	if err != nil {
		bye(fmt.Sprintf("could not build site: %s", err.Error()), 1)
	}
}
