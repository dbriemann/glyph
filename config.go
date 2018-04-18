package main

type Repository struct {
	Users []string `toml:"Users"` // first must be login name
	Name  string   `toml:"Name"`
}

type Site struct {
	Title       string `toml:"Title"`
	Author      string `toml:"Author"`
	OneLineDesc string `toml:"OneLineDesc"`
	Twitter     string `toml:"Twitter"`
	Mail        string `toml:"Mail"`
	Theme       string `toml:"Theme"`
}

type Template struct {
	Source string `toml:"Source"`
	Layout string `toml:"Layout"`
	Target string `toml:"Target"`
}

type ThemeConfig struct {
	Name           string                 `toml:"Name"`
	IndexTemplate  Template               `toml:"IndexTemplate"`
	IssueTemplate  Template               `toml:"IssueTemplate"`
	OtherTemplates []Template             `toml:"OtherTemplates"`
	Custom         map[string]interface{} `toml:"Custom"`
}

type BaseConfig struct {
	Repository  Repository             `toml:"Repository"`
	Site        Site                   `toml:"Site"`
	Custom      map[string]interface{} `toml:"Custom"`
	GithubToken string                 `toml:"GithubToken"`
}
