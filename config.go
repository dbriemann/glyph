package main

type Repository struct {
	Users []string `toml:"users"` // first must be login name
	Name  string   `toml:"name"`
}

type Site struct {
	Title       string `toml:"title"`
	Author      string `toml:"author"`
	OneLineDesc string `toml:"onelinedesc"`
	Twitter     string `toml:"twitter"`
	Mail        string `toml:"mail"`
	Theme       string `toml:"theme"`
}

type Config struct {
	Repository Repository             `toml:"repository"`
	Site       Site                   `toml:"site"`
	Custom     map[string]interface{} `toml:"custom"`
}
