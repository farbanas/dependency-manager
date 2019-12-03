package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var username string
var password string
var nexusUrl string
var configPath string

func init() {
	username = os.Getenv("USERNAME")
	password = os.Getenv("PASSWORD")
	nexusUrl = os.Getenv("NEXUS_URL")
	configPath = os.Getenv("CONFIG_PATH")

	flag.StringVar(&username, "username", username, "Username for the Nexus repository")
	flag.StringVar(&password, "password", password, "Password for the Nexus repository")
	flag.StringVar(&nexusUrl, "url", nexusUrl, "URL for the Nexus repository")
	flag.StringVar(&configPath, "config-path", configPath, "Path to config.yaml")
	flag.Parse()
}

func parseVars() {
	if username == "" {
		panic("Missing username for nexus repo!")
	}
	if password == "" {
		panic("Missing password for nexus repo!")
	}
	if nexusUrl == "" {
		panic("Missing nexus url!")
	}
}

func main() {
	parseVars()
	config := ReadConfigFile("config.yaml")
	for i, repoConfig := range config.Config {
		if i != 0 {
			fmt.Println()
		}
		fmt.Println(strings.Title(repoConfig.RepositoryType), "packages that need update:")
		strlen := len(repoConfig.RepositoryType) + len("packages that need update:")
		fmt.Println(strings.Repeat("=", strlen+1))
		for _, reqFileDef := range repoConfig.RequirementFiles {
			var repository Requirements
			if repoConfig.RepositoryType == "helm" {
				repository = HelmRequirements{reqFileDef.Path, nil, HelmYaml{}}
			} else if repoConfig.RepositoryType == "pip" {
				repository = PipRequirements{reqFileDef.Path, nil }
			} else if repoConfig.RepositoryType == "pom" {
				repository = PomRequirements{reqFileDef.Path, nil, PomXML{}}
			}
			if repository != nil {
				process(reqFileDef, repository)
			} else {
				fmt.Println("No packages need update!")
			}
		}
	}
}

func process(reqDefinition RequirementsDefinition, requirements Requirements) {
	toUpdate := make(map[string]string)
	reader, f := requirements.OpenRequirementsFile()
	defer f.Close()
	requirements = requirements.ReadCurrentVersion(reader)

	libraryVersions := requirements.GetLibraryVersions()
	for _, req := range libraryVersions {
		url := fmt.Sprintf("https://%s:%s@%s/service/rest/v1/search?repository=%s&name=%s&sort=version", username, password, nexusUrl, reqDefinition.Repository, req.Library)
		assets := NexusGetAssets(url)
		versions := NexusExtractVersions(assets)

		if len(versions) > 0 {
			if versions[0] != req.Version {
				toUpdate[req.Library] = versions[0]
				fmt.Printf("%-25s\t%-7s -> %7s\n", req.Library, req.Version, versions[0])
			}
		} else {
			fmt.Println("Could not get version for library:", req.Library)
		}
	}
	fmt.Println("Getting done, starting update")
	requirements.UpdateVersion(toUpdate)
}
