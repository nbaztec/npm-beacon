package handler

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration class
type Configuration struct {
	GithubToken       string   `yaml:"githubToken"`
	Repositories      []string `yaml:"repositories"`
	MinDaysNewRelease int      `yaml:"minDaysNewRelease"`
}

// LoadConfiguration reads the configuration file
func LoadConfiguration() *Configuration {
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
