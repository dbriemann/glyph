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

type Theme struct {
	Name string `toml:"Name"`
}

type Config struct {
	Repository Repository             `toml:"Repository"`
	Site       Site                   `toml:"Site"`
	Custom     map[string]interface{} `toml:"Custom"`
}
