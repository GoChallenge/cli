package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/cli"
)

func fetch(c *cli.Context) {
	currentChallengeApiEndpoint := strings.Join([]string{apiUrl, "current"}, "/")
	if err := getChallenge(currentChallengeApiEndpoint); err != nil {
		fmt.Println(err.Error())
	}

	return
}

func getChallenge(apiPath string) error {
	apiResponse, err := http.Get(apiPath)
	defer apiResponse.Body.Close()
	if err != nil {
		return fmt.Errorf("Unable to contact API. Error: %s\n", err.Error())
	}

	body, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		return fmt.Errorf("Api responded with an error. Status code: %d\n", apiResponse.StatusCode)
	}

	if apiResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("Api responded with an error. Status code: %d. Body: %s\n", apiResponse.StatusCode, string(body))
	}

	var challengeData challengeDataStruct
	if err := json.Unmarshal(body, &challengeData); err != nil {
		return fmt.Errorf("Error while reading response from api: %s\n", err.Error())
	}

	if output, err := downloadChallenge(challengeData.Import); err != nil {
		return fmt.Errorf("Unable to `go get` challenge from import path %s. Command output: %s\n", challengeData.Import, output)
	}

	dir, err := getPackageDirectory(challengeData.Import)
	if err != nil {
		return fmt.Errorf("Error determining downloaded directory: %s\n", err.Error())
	}

	fmt.Printf("Downloaded the latest challenge to %s. See README.md inside the directory or go to %s for information on the challenge\n", dir, challengeData.Url)

	return nil
}

func downloadChallenge(importPath string) (string, error) {
	goGetCmd := exec.Command("go", "get", importPath)

	output, err := goGetCmd.CombinedOutput()
	return string(output), err
}

func getPackageDirectory(importPath string) (string, error) {
	pkg, err := build.Import(importPath, os.Getenv("GOPATH"), build.FindOnly)
	if err != nil {
		return "", err
	}

	return pkg.Dir, nil
}
