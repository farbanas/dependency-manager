package main

import (
	json2 "encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type nexusItems struct {
	Items             []nexusComponent `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

type nexusComponent struct {
	Id         string       `json:"id"`
	Repository string       `json:"repository"`
	Format     string       `json:"format"`
	Group      string       `json:"group"`
	Name       string       `json:"name"`
	Version    string       `json:"version"`
	Assets     []nexusAsset `json:"assets"`
}

type nexusAsset struct {
	DownloadUrl string        `json:"download_url"`
	Path        string        `json:"path"`
	Id          string        `json:"id"`
	Repository  string        `json:"repository"`
	Format      string        `json:"format"`
	Checksum    nexusChecksum `json:"checksum"`
}

type nexusChecksum struct {
	Sha1   string `json:"sha_1"`
	Sha256 string `json:"sha_256"`
	Md5    string `json:"md_5"`
}

func NexusGetAssets(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	components, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return components
}

func NexusExtractVersions(components []byte) []string {
	var nComponents nexusItems
	err := json2.Unmarshal(components, &nComponents)
	if err != nil {
		log.Fatal(err)
	}

	versions := make([]string, 0, cap(nComponents.Items))
	for _, item := range nComponents.Items {
		versions = append(versions, item.Version)
	}
	return versions
}

func sortSemVer(versions []string, i, j int, reverse bool) bool {
	//TODO: this function should only sort if versions are true semVer
	splitVersions1 := strings.Split(versions[i], ".")
	splitVersions2 := strings.Split(versions[j], ".")
	major1, _ := strconv.Atoi(string(splitVersions1[0]))
	major2, _ := strconv.Atoi(string(splitVersions2[0]))
	if reverse {
		if major1 > major2 {
			return true
		}
	} else {
		if major1 < major2 {
			return true
		}
	}
	if len(splitVersions1) >= 2 && len(splitVersions2) >= 2 {
		minor1, _ := strconv.Atoi(string(splitVersions1[1]))
		minor2, _ := strconv.Atoi(string(splitVersions2[1]))
		if reverse {
			if minor1 > minor2 {
				return true
			}
		} else {
			if minor1 < minor2 {
				return true
			}
		}
	}

	if len(splitVersions1) == 3 && len(splitVersions2) == 3 {
		patch1, _ := strconv.Atoi(string(splitVersions1[2]))
		patch2, _ := strconv.Atoi(string(splitVersions2[2]))
		if reverse {
			if patch1 > patch2 {
				return true
			}
		} else {
			if patch1 < patch2 {
				return true
			}
		}
	}
	return false
}
