package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/codegangsta/cli"
)

const (
	loginFileName = ".gochallenge.json"
	defaultHost   = ""
)

type loginDetails struct {
	GithubUsername string `json:"gihubUsername"`
	ApiKey         string `json:"apiKey"`
	Hostname       string `json:"hostname"`
}

func login(c *cli.Context) {
	username, err := getString("Your GitHub username:")
	if err != nil {
		fmt.Println(err)
		return
	}
	apikey, err := getString("Your GoChallenge API key (found at http://golang-challenge.com/account):")
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO
	// check with server if this is ok

	ld := loginDetails{username, apikey, defaultHost}
	lds, err := json.Marshal(ld)
	if err != nil {
		fmt.Println(err)
		return
	}

	loginfile, err := getLoginFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(loginfile, lds, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully logged in")
}

func logout(c *cli.Context) {
	loginfile, err := getLoginFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.Remove(loginfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully logged out")
}

func getLoginDetails() (loginDetails, error) {
	loginfile, err := getLoginFile()
	if err != nil {
		return loginDetails{}, err
	}
	data, err := ioutil.ReadFile(loginfile)
	if err != nil {
		return loginDetails{}, err
	}
	var details loginDetails
	err = json.Unmarshal(data, &details)
	return details, err
}

func getLoginFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, loginFileName), nil
}

func getString(prompt string) (string, error) {
	fmt.Print(prompt)
	var result string
	_, err := fmt.Scan(&result)
	return result, err
}
