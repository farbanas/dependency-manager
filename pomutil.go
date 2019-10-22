package main

import (
	"bufio"
	"os"
)

type Pom struct {
	path            string
	libraryVersions []LibraryVersion
}

func (p Pom) OpenPomFile() *bufio.Reader {
	f, err := os.Open(p.path)
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(f)
	return reader
}

func (p Pom) ReadCurrentVersion(reader *bufio.Reader) {
}
