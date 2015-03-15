package cmd

import (
	"encoding/json"
	"errors"
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
	Name   string `json:"name"`
	URL    string `json:"url"`
	Status string `json:"status"`
	Import string `json:"import"`
}

// Fetch fetches the details of the current challenge and stores it in a file (challengeFile)
// in the home directory
func Fetch(c *cli.Context) {
	if c.Bool("older") {
		err := fetchOlder()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	chal, err := fetchLatest()
	if err != nil {
		fmt.Println(err)
		return
	}

	if output, err := chal.download(); err != nil {
		fmt.Printf("Unable to `go get` challenge %s\n", chal.Import)
		fmt.Println(output)
		return
	}
	fmt.Printf("Downloaded the latest challenge to %s\n", chal.directory())
	fmt.Printf("See README.md or go to %s for information\n", chal.URL)

	if err = chal.store(); err != nil {
		fmt.Printf("Could not store challenge data - %s\n", err)
	}
}

func fetchOlder() error {
	challenges, err := fetchChallenges(apiURL)
	if err != nil {
		return err
	}
	if len(challenges) == 0 {
		return errors.New("Found no challenges")
	}
	for _, chal := range challenges {
		fmt.Printf("Downloading %s\n", chal.Name)
		if output, err := chal.download(); err != nil {
			fmt.Printf("Unable to `go get` challenge %s\n", chal.Import)
			fmt.Println(output)
		}
	}
	return nil
}

func fetchLatest() (challenge, error) {
	currentURL := strings.Join([]string{apiURL, "current"}, "/")
	chals, err := fetchChallenges(currentURL)
	if err != nil {
		return challenge{}, err
	}
	return chals[0], nil
}

func fetchChallenges(url string) ([]challenge, error) {
	var chal []challenge
	resp, err := http.Get(url)
	if err != nil {
		return chal, fmt.Errorf("Unable to contact API. Error: %s\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return chal, err
	}

	if resp.StatusCode != http.StatusOK {
		return chal, fmt.Errorf("Api responded with an error. Status code: %d. Body: %s\n", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &chal); err != nil {
		return chal, fmt.Errorf("Error while reading response from api: %s\n", err)
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
