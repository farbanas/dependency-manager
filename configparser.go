package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Config []RepositoryConfig `yaml:"config"`
}

type RepositoryConfig struct {
	RepositoryType   string                   `yaml:"type"`
	RequirementFiles []RequirementsDefinition `yaml:"spec"`
}

type RequirementsDefinition struct {
	Repository string `yaml:"repository"`
	Path       string `yaml:"path"`
	AutoUpdate string `yaml:"auto-update"`
}

func ReadConfigFile(path string) Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}
