package main

import (
	"bufio"
	"os"
)

type Requirements interface {
	ReadCurrentVersion(reader *bufio.Reader) Requirements
	OpenRequirementsFile() (*bufio.Reader, *os.File)
	GetLibraryVersions() []LibraryVersion
	GetPath() string
	UpdateVersion(map[string]string)
}

type FileData interface {
	WriteData()
	GetDependencies()
}

type LibraryVersion struct {
	Library string
	Version string
}
