package main

import (
	"fmt"
	"os"
)

var USERNAME string
var PASSWORD string
var NEXUSURL string

func main() {
	USERNAME = os.Getenv("USERNAME")
	PASSWORD = os.Getenv("PASSWORD")
	NEXUSURL = os.Getenv("NEXUS_URL")

	config := ReadConfigFile("config.yaml")

	for _, repoConfig := range config.Config {
		if repoConfig.RepositoryType == "helm" {
			for _, reqFileDef := range repoConfig.RequirementFiles {
				repository := HelmRequirements{reqFileDef.Path, nil}
				process(reqFileDef, repository)
			}
		} else if repoConfig.RepositoryType == "pip" {
			for _, reqFileDef := range repoConfig.RequirementFiles {
				repository := PipRequirements{reqFileDef.Path, nil}
				process(reqFileDef, repository)
			}
		}
	}
}

func process(reqDefinition RequirementsDefinition, requirements Requirements) {
	reader, f := requirements.OpenRequirementsFile()
	defer f.Close()
	requirements = requirements.ReadCurrentVersion(reader)

	libraryVersions := requirements.GetLibraryVersions()
	for _, req := range libraryVersions {
		url := fmt.Sprintf("https://%s:%s@%s/service/rest/v1/search?repository=%s&name=%s&sort=version", USERNAME, PASSWORD, NEXUSURL, reqDefinition.Repository, req.Library)
		assets := NexusGetAssets(url)
		versions := NexusExtractVersions(assets)

		if len(versions) > 0 {
			if versions[0] != req.Version {
				fmt.Printf("%-25s\t%-7s -> %7s\n", req.Library, req.Version, versions[0])
			}
		} else {
			fmt.Println("Could not get version for library:", req.Library)
		}
	}
}
