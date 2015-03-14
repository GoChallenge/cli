package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/codegangsta/cli"
)

type challenge struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Status string `json:"status"`
	Import string `json:"import"`
}

func fetch(c *cli.Context) {
	currentURL := strings.Join([]string{apiUrl, "current"}, "/")
	chal, err := getChallenge(currentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if output, err := chal.download(); err != nil {
		fmt.Println(fmt.Sprintf("Unable to `go get` challenge %s", chal.Import))
		fmt.Println(output)
		return
	} else {
		fmt.Println(fmt.Sprintf("Downloaded the latest challenge to %s", chal.directory()))
		fmt.Println(fmt.Sprintf("See README.md inside the directory or go to %s for information on the challenge", chal.Url))
	}

	chal.store()
}

func getChallenge(url string) (challenge, error) {
	var chal challenge
	resp, err := http.Get(url)
	if err != nil {
		return chal, fmt.Errorf("Unable to contact API. Error: %s\n", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return chal, fmt.Errorf("Api responded with an error. Status code: %d\n", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return chal, fmt.Errorf("Api responded with an error. Status code: %d. Body: %s\n", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &chal); err != nil {
		return chal, fmt.Errorf("Error while reading response from api: %s\n", err.Error())
	}
	return chal, nil
}

func (c challenge) download() ([]byte, error) {
	goGetCmd := exec.Command("go", "get", c.Import)
	return goGetCmd.CombinedOutput()
}

func (c challenge) directory() string {
	return path.Join(os.Getenv("GOPATH"), "src", c.Import)
}

func (c challenge) store() error {
	return nil
}
