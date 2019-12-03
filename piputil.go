package main

import (
	"bufio"
	"os"
	"strings"
)

type PipRequirements struct {
	path            string
	libraryVersions []LibraryVersion
}

func (p PipRequirements) OpenRequirementsFile() (*bufio.Reader, *os.File) {
	f, err := os.Open(p.path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	return reader, f
}

func (p PipRequirements) ReadCurrentVersion(reader *bufio.Reader) Requirements {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var tokens []string
		line := scanner.Text()
		if tokens = tokenize(line); tokens != nil {
			p.libraryVersions = append(p.libraryVersions, LibraryVersion{tokens[0], tokens[1]})
		}
	}
	return p
}

func (p PipRequirements) GetLibraryVersions() []LibraryVersion {
	return p.libraryVersions
}

func (p PipRequirements) GetPath() string {
	return p.path
}

func tokenize(line string) []string {
	line = strings.Trim(line, " ")
	var tokens []string
	if line == "" {
		tokens = nil
	} else if strings.HasPrefix(line, "#") {
		tokens = nil
	} else if strings.HasPrefix(line, "--") {
		tokens = nil
	} else if strings.Contains(line, "==") {
		tokens = strings.Split(line, "==")
	} else if strings.Contains(line, ">=") {
		tokens = strings.Split(line, ">=")
	} else if strings.Contains(line, "<=") {
		tokens = strings.Split(line, "<=")
	} else if strings.Contains(line, "<") {
		tokens = strings.Split(line, "<")
	} else if strings.Contains(line, ">") {
		tokens = strings.Split(line, ">")
	} else {
		tokens = []string{line, "?"}
	}
	return tokens
}

func (p PipRequirements) UpdateVersion(toUpdate map[string]string){

}

