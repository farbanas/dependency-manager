package main

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type HelmDependency struct {
	Name       string `yaml:"name"`
	Repository string `yaml:"repository"`
	Version    string `yaml:"version"`
}

type HelmYaml struct {
	Dependencies []HelmDependency `yaml:"dependencies"`
}

type HelmRequirements struct {
	path            string
	libraryVersions []LibraryVersion
	data 			HelmYaml
}

func (h HelmRequirements) OpenRequirementsFile() (*bufio.Reader, *os.File) {
	f, err := os.Open(h.path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	return reader, f
}

func (h HelmRequirements) ReadCurrentVersion(reader *bufio.Reader) Requirements {
	var dependencies HelmYaml
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &dependencies)
	if err != nil {
		panic(err)
	}

	h.data = dependencies

	return h.UnpackDependencies(dependencies)
}

func (h HelmRequirements) UnpackDependencies(dependencies HelmYaml) HelmRequirements {
	for _, dependency := range dependencies.Dependencies {
		if dependency.Repository != "@stable" {
			h.libraryVersions = append(h.libraryVersions, LibraryVersion{dependency.Name, dependency.Version})
		}
	}
	return h
}

func (h HelmRequirements) UpdateDependencies(dependencies HelmYaml) HelmRequirements {
	for _, dependency := range dependencies.Dependencies {
		if dependency.Repository != "@stable" {
			h.libraryVersions = append(h.libraryVersions, LibraryVersion{dependency.Name, dependency.Version})
		}
	}
	return h
}

func (h HelmRequirements) GetLibraryVersions() []LibraryVersion {
	return h.libraryVersions
}

func (h HelmRequirements) GetPath() string {
	return h.path
}

func (h HelmRequirements) UpdateVersion(toUpdate map[string]string) {
	for i, dependency := range h.data.Dependencies {
		if version := toUpdate[dependency.Name]; version != "" {
			h.data.Dependencies[i].Version = version
		}
	}
	yamlData, err := yaml.Marshal(h.data)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(h.path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(yamlData)
	if err != nil {
		panic(err)
	}
}