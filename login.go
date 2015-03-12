package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
)

const (
	logindetailsFile = "/Users/krishnasundarram/.gochallenge.go"
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

	ld := loginDetails{username, apikey, ""}
	lds, err := json.Marshal(ld)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(logindetailsFile, lds, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully logged in")
}

func logout(c *cli.Context) {
	err := os.Remove(logindetailsFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully logged out")
}

func getLoginDetails() (loginDetails, error) {
	data, err := ioutil.ReadFile(logindetailsFile)
	if err != nil {
		return loginDetails{}, err
	}
	var details loginDetails
	err = json.Unmarshal(data, &details)
	return details, err
}

func getString(prompt string) (string, error) {
	fmt.Print(prompt)
	var result string
	_, err := fmt.Scan(&result)
	return result, err
}
