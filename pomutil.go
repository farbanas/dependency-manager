package main

import (
	"bufio"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
)

type PomDependency struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type PomXML struct {
		Dependencies []PomDependency `xml:"dependencies>dependency"`
}

type PomRequirements struct {
	path            string
	libraryVersions []LibraryVersion
}

func (p PomRequirements) OpenRequirementsFile() (*bufio.Reader, *os.File) {
	f, err := os.Open(p.path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	return reader, f
}

func (p PomRequirements) ReadCurrentVersion(reader *bufio.Reader) Requirements {
	var dependencies PomXML
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	err = xml.Unmarshal(data, &dependencies)
	return p.UnpackDependencies(dependencies)
}

func (p PomRequirements) UnpackDependencies (dependencies PomXML) PomRequirements {
	for _, dependency := range dependencies.Dependencies {
		if strings.Contains(dependency.GroupId, "com.vingd") {
			p.libraryVersions = append(p.libraryVersions, LibraryVersion {
				Library: dependency.ArtifactId,
				Version: dependency.Version,
			})
		}
	}
	return p
}

func (p PomRequirements) GetLibraryVersions() []LibraryVersion {
	return p.libraryVersions
}

func (p PomRequirements) GetPath() string {
	return p.path
}
