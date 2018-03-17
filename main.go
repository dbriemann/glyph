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
	baseDir    = "includes"
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
// TODO 4. add syntax highlighting with chroma

func main() {
	// Read config file and check sanity (TODO).
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err.Error())
	}
	if _, err := toml.Decode(string(raw), &cfg); err != nil {
		panic(err.Error())
	}

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
	} else {
		if !finfo.IsDir() {
			panic(fmt.Sprintf("%s should be a directory but is a file", cfg.Repository.TargetDir))
		}
	}

	// Copy include files (css and js stuff).
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		panic(err.Error())
	}

	for _, f := range files {
		src := filepath.Join(baseDir, f.Name())
		dst := filepath.Join(cfg.Repository.TargetDir, f.Name())
		if err := copyFile(src, dst); err != nil {
			panic(err.Error())
		}
	}

	// We don't use access tokens because the rate limiting for unauthed access is good enough.
	// This way we have an easy time using this in CI scripts without having to provide secret
	// information.
	ctx := context.Background()
	client := github.NewClient(nil)

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
