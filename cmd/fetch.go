package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	var err error
	if c.Bool("all") {
		err = fetchAll()
	} else if c.IsSet("challenge") {
		_, err = fetchChallenge(c.String("challenge"))
	} else {
		err = fetchCurrent()
	}
	if err != nil {
		fmt.Println(err)
	}
}

func fetchAll() error {
	challenges, err := getChallengeDescriptors(apiURL)
	if err != nil {
		return err
	}
	return fetchChallenges(challenges)
}

func fetchCurrent() error {
	chal, err := fetchChallenge("current")
	if err != nil {
		return err
	}
	fmt.Printf("See README.md or go to %s for information\n", chal.URL)

	if err = chal.store(); err != nil {
		return fmt.Errorf("Could not store challenge data - %s\n", err)
	}
	return nil
}

func fetchChallenge(ID string) (*challenge, error) {
	challenges, err := getChallengeDescriptors(strings.Join([]string{apiURL, ID}, "/"))
	if err != nil {
		return nil, err
	}
	if len(challenges) != 1 {
		return nil, fmt.Errorf("Found multiple challenges in response to single challenge query for challenge \"%s\"", ID)
	}
	err = fetchChallenges(challenges)
	return &challenges[0], nil
}

func fetchChallenges(descriptors []challenge) error {
	for _, chal := range descriptors {
		fmt.Printf("Downloading %s to %s ........ ", chal.Name, chal.directory())
		if output, err := chal.download(); err != nil {
			fmt.Println()
			return fmt.Errorf("Unable to `go get` challenge %s\n%s", chal.Import, output)
		}
		fmt.Println("Done")
	}
	return nil
}

func getChallengeDescriptors(url string) ([]challenge, error) {
	var challenges []challenge
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to contact API. Error: %s\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 100*1024)) // limit response size to 100KB
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Api responded with an error. Status code: %d\n", resp.StatusCode)
	}
	c, err := decodeSingleDescriptor(body)
	if err != nil {
		// Couldn't decode as a single challenge descriptor, try as an array.
		challenges, err = decodeMultipleDescriptors(body)
		if err != nil {
			return nil, err
		}
	} else {
		challenges = append(challenges, *c)
	}

	if len(challenges) == 0 {
		return nil, errors.New("Found no challenges")
	}

	return challenges, nil
}

func decodeSingleDescriptor(body []byte) (*challenge, error) {
	var c challenge
	err := json.Unmarshal(body, &c)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response from api: %s\n", err)
	}
	return &c, nil
}

func decodeMultipleDescriptors(body []byte) ([]challenge, error) {
	var challenges []challenge
	err := json.Unmarshal(body, &challenges)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response from api: %s\n", err)
	}
	return challenges, nil
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
