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
var configPath	 string


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

func main() {
	config := ReadConfigFile("config.yaml")
	for i, repoConfig := range config.Config {
		if i != 0 {
			fmt.Println()
		}
		fmt.Println(strings.Title(repoConfig.RepositoryType), "packages that need update:")
		strlen := len(repoConfig.RepositoryType) + len("packages that need update:")
		fmt.Println(strings.Repeat("=", strlen + 1))


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
		} else if repoConfig.RepositoryType == "pom" {
			for _, reqFileDef := range repoConfig.RequirementFiles {
				repository := PomRequirements{reqFileDef.Path, nil}
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
		url := fmt.Sprintf("https://%s:%s@%s/service/rest/v1/search?repository=%s&name=%s&sort=version", username, password, nexusUrl, reqDefinition.Repository, req.Library)
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
