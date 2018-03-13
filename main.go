package main

import (
	"context"
	"fmt"

	"os"

	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	configFile = "config.toml"
)

var (
	cfg    Config
	outDir string
)

// TODO 1. replace explicit panics with error messages and proper handling
// TODO 2. add cli options
// TODO 2.a. option: init project
// TODO 2.b. option: build project
// TODO 3. check config data for sanity
// TODO 4. slugify issue titles

func main() {
	// Read config file and check sanity (TODO).
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err.Error())
	}
	if _, err := toml.Decode(string(raw), &cfg); err != nil {
		panic(err.Error())
	}
	// Fetch Github access token from environment.
	token := os.Getenv(cfg.GithubToken)

	// Set output directory.
	finfo, err := os.Stat(cfg.Repository.TargetDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(cfg.Repository.TargetDir, 0755)
			if err != nil {
				panic(err.Error())
			}
		} else {
			panic(err.Error())
		}
	}
	if finfo == nil || !finfo.IsDir() {
		panic(fmt.Sprintf("%s should be a directory but is a file", cfg.Repository.TargetDir))
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if client == nil {
		panic("client not working")
	}

	issues, err := FetchIssues(client, ctx, &cfg)
	if err != nil {
		panic(err.Error())
	}

	err = BuildSite(issues, &cfg)
	if err != nil {
		panic(err.Error())
	}
}
