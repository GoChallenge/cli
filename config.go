package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	configFileName = ".gochallenge.json"
)

type configData struct {
	ApiKey string `json:"apiKey"`
}

func (cd *configData) printData() {
	fmt.Printf("Api Key: %s\n", cd.ApiKey)
}

func showConfig() {
	currentConfigData, err := getCurrentConfig()
	if err != nil {
		fmt.Printf("Error while reading current config: %s\n", err.Error())
		return
	}

	currentConfigData.printData()
}

func writeConfig(c *cli.Context) {
	apiKey := c.String("apiKey")
	if len(apiKey) == 0 {
		showConfig()
		return
	}

	newConfigData := &configData{
		ApiKey: apiKey,
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return
	}
	_, err = os.Stat(configFilePath)
	// We need to confirm before we overwrite an existing file
	if err == nil {
		fmt.Printf("The file %s already exists. Are you sure you want to overwride it? (yes/no)\n", configFilePath)

		var resp string
		_, err = fmt.Scanf("%s", &resp)
		if err != nil {
			fmt.Printf("Error while reading user input: %s\n", err.Error())
			return
		}

		resp = strings.TrimSpace(resp)
		if resp != "yes" {
			fmt.Println("Exiting without writing to config file")
			return
		}
	} else if !os.IsNotExist(err) {
		fmt.Printf("Unable to get info on file %s: %s\n", configFilePath, err.Error())
		return
	}

	configFile, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Unable to open file %s for writing: %s\n", configFilePath, err.Error())
		return
	}

	jsonString, err := json.Marshal(newConfigData)
	if err != nil {
		fmt.Printf("Unable to convert config data to json: %s\n", err.Error())
		return
	}

	_, err = configFile.Write(jsonString)
	if err != nil {
		fmt.Printf("Error while writing to config file %s: %s\n", configFilePath, err.Error())
		return
	}

	err = configFile.Close()
	if err != nil {
		fmt.Printf("Unable to close config file %s: %s\n", configFileName, err.Error())
		return
	}

	fmt.Printf("Config written to %s. You can now use the rest of the features of the CLI app. Happy hacking!\n", configFilePath)
}

func getConfigFilePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Unable to get current user. Error: %s\n", err.Error())
		return "", err
	}

	userHomeDir := currentUser.HomeDir
	configFilePath := path.Join(userHomeDir, configFileName)

	return configFilePath, nil
}

func getCurrentConfig() (*configData, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	configFileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("It seems you haven't create a config file yet. Please create one by using the `config` command")
		}
		return nil, err
	}

	currentConfigData := new(configData)
	err = json.Unmarshal(configFileContent, currentConfigData)
	if err != nil {
		return nil, err
	}

	return currentConfigData, nil
}
