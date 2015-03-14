package cmd

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
	configFile = ".gochallenge_config.json" // Will be stored in the home directory
)

type config struct {
	APIKey string `json:"apiKey"`
}

// Configure stores the APIKey in a file (configFile) in the home directory
func Configure(c *cli.Context) {
	apikey := c.String("key")
	if apikey == "" {
		fmt.Println("Empty key")
		return
	}

	cfg := config{apikey}
	cfgs, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfgfile, err := getConfigFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(cfgfile, cfgs, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully configured")
}

func readConfigFile() (config, error) {
	var cfg config

	cfgfile, err := getConfigFile()
	if err != nil {
		return cfg, err
	}
	data, err := ioutil.ReadFile(cfgfile)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func getConfigFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, configFile), nil
}
