package main

type Repository struct {
	Users     []string `toml:"users"`
	Name      string   `toml:"name"`
	TargetDir string   `toml:"target_dir"`
}

type Site struct {
	Title       string `toml:"title"`
	Author      string `toml:"author"`
	OneLineDesc string `toml:"onelinedesc"`
	Twitter     string `toml:"twitter"`
}

type Config struct {
	GithubToken string     `toml:"github_token"`
	Repository  Repository `toml:"repository"`
	Site        Site       `toml:"site"`
}
