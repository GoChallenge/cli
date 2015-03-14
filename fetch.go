package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	challengeFile = ".gochallenge_data.json" // Will be stored in the home directory
)

type challenge struct {
	ID     int    `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"`
	Import string `json:"import"`
}

func fetch(c *cli.Context) {
	currentURL := strings.Join([]string{apiURL, "current"}, "/")
	fmt.Println(currentURL)
	chal, err := fetchChallenge(currentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if output, err := chal.download(); err != nil {
		fmt.Println(fmt.Sprintf("Unable to `go get` challenge %s", chal.Import))
		fmt.Println(output)
		return
	}
	fmt.Println(fmt.Sprintf("Downloaded the latest challenge to %s", chal.directory()))
	fmt.Println(fmt.Sprintf("See README.md inside the directory or go to %s for information on the challenge", chal.URL))

	if err = chal.store(); err != nil {
		fmt.Println(fmt.Sprintf("Could not store challenge data - %s", err.Error()))
	}
}

func fetchChallenge(url string) (challenge, error) {
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

func readChallengeFile() (challenge, error) {
	var chal challenge

	filepath, err := getChallengeFile()
	if err != nil {
		return chal, err
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return chal, err
	}
	err = json.Unmarshal(data, &chal)
	return chal, err
}

func getChallengeFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, challengeFile), nil
}

func (c challenge) download() ([]byte, error) {
	goGetCmd := exec.Command("go", "get", c.Import)
	return goGetCmd.CombinedOutput()
}

func (c challenge) directory() string {
	return path.Join(os.Getenv("GOPATH"), "src", c.Import)
}

func (c challenge) store() error {
	filepath, err := getChallengeFile()
	if err != nil {
		return err
	}

	chal, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, chal, os.ModePerm)
}
