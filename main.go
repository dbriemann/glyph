package main

import (
	"context"
	"fmt"
	"path/filepath"

	"os"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
)

const (
	feedFile   = "feed.atom"
	baseDir    = "themes"
	configFile = "config.toml"
)

var (
	cfg      Config
	themeDir = ""
	outDir   string
)

// TODO add syntax highlighting with chroma?

func main() {
	// Read config file and check sanity (TODO).
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		bye(fmt.Sprintf("could not read config file: %s", err.Error()), 1)
	}
	if _, err := toml.Decode(string(raw), &cfg); err != nil {
		bye(fmt.Sprintf("could not parse config file: %s", err.Error()), 1)
	}

	// Test config data sanity.
	if cfg.Repository.Name == "" {
		bye("config file: repository name missing", 1)
	}
	if cfg.Repository.OutputDir == "" {
		bye("config file: output directory missing", 1)
	}
	if len(cfg.Repository.Users) < 1 {
		bye("config file: no user(s) provided", 1)
	}
	if cfg.Site.Title == "" {
		fmt.Println("warning: no site title set")
	}
	if cfg.Site.Author == "" {
		fmt.Println("warning: no site author set")
	}
	if cfg.Site.OneLineDesc == "" {
		fmt.Println("warning: no one line description set")
	}
	if cfg.Site.Mail == "" {
		fmt.Println("warning: no contact mail set")
	}
	if cfg.Site.Twitter == "" {
		fmt.Println("warning: no twitter handle set")
	}
	if cfg.Site.Theme == "" {
		fmt.Println("warning: no theme set.. falling back to default")
		cfg.Site.Theme = "default"
	}

	// Set output directory.
	finfo, err := os.Stat(cfg.Repository.OutputDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(cfg.Repository.OutputDir, 0755)
			if err != nil {
				bye(fmt.Sprintf("could not create output directory: %s", err.Error()), 1)
			}
		} else {
			bye(fmt.Sprintf("could not access output directory: %s", err.Error()), 1)
		}
	} else {
		if !finfo.IsDir() {
			bye(fmt.Sprintf("%s should be a directory but is a file", cfg.Repository.OutputDir), 1)
		}
	}

	// Try to open theme folder..
	themeDir = filepath.Join(baseDir, cfg.Site.Theme)
	finfo, err = os.Stat(themeDir)
	if err != nil || !finfo.IsDir() {
		fmt.Sprintf("cannot load theme %s: %s\n", cfg.Site.Theme, err.Error())
		fmt.Println("falling back to default theme")
	}
	themeDir = filepath.Join(baseDir, "default")

	// Copy theme files except mustache template files.
	files, err := ioutil.ReadDir(themeDir)
	if err != nil {
		bye(fmt.Sprintf("could not read access directory: %s", err.Error()), 1)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".mustache" {
			src := filepath.Join(themeDir, f.Name())
			dst := filepath.Join(cfg.Repository.OutputDir, f.Name())
			if err := copyFile(src, dst); err != nil {
				bye(fmt.Sprintf("could not copy file: %s", err.Error()), 1)
			}
		}
	}

	// We don't use access tokens because the rate limiting for unauthed access is good enough.
	// This way we have an easy time using this in CI scripts without having to provide secret
	// information.
	ctx := context.Background()
	client := github.NewClient(nil)

	if client == nil {
		bye("client dead", 1)
	}

	issues, err := FetchIssues(client, ctx, cfg)
	if err != nil {
		bye(fmt.Sprintf("could not fetch issues from github: %s", err.Error()), 1)
	}

	err = BuildSite(issues, cfg)
	if err != nil {
		bye(fmt.Sprintf("could not build site: %s", err.Error()), 1)
	}
}
