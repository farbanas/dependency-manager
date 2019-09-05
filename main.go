package main

import (
	json2 "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type nexusItems struct {
	Items	[]nexusComponent	`json:"items"`
	ContinuationToken	string	`json:"continuationToken"`
}

type nexusComponent struct {
	Id         string	`json:"id"`
	Repository string	`json:"repository"`
	Format     string	`json:"format"`
	Group      string	`json:"group"`
	Name       string	`json:"name"`
	Version    string	`json:"version"`
	Assets     []nexusAsset	`json:"assets"`
}

type nexusAsset struct {
	DownloadUrl string	`json:"download_url"`
	Path        string	`json:"path"`
	Id          string	`json:"id"`
	Repository  string	`json:"repository"`
	Format      string	`json:"format"`
	Checksum    nexusChecksum	`json:"checksum"`
}

type nexusChecksum struct {
	Sha1   string	`json:"sha_1"`
	Sha256 string	`json:"sha_256"`
	Md5    string	`json:"md_5"`
}

func main() {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	repo := os.Getenv("REPO")
	asset := os.Getenv("ASSET")
	url := fmt.Sprintf("https://%s:%s@nexus.vingd.net/service/rest/v1/search?repository=%s&name=%s&sort=version", username, password, repo, asset)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	components, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	nexusExtractVersions(components)
}

func nexusGetPackages() {
}

func nexusExtractVersions(components []byte) {
	var nComponents nexusItems
	err := json2.Unmarshal(components, &nComponents)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nComponents.Items[0].Version)
}
