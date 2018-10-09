package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration class
type Configuration struct {
	GithubToken  string   `yaml:"githubToken"`
	Repositories []string `yaml:"repositories"`
}

// Load configuration file
func Load() *Configuration {
	filepath := "./config.yaml"
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var config Configuration
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return &config
}
