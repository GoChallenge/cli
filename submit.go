package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

func submit(c *cli.Context) {
	config, err := getConfig()
	if err != nil || config.ApiKey == "" {
		fmt.Println("Please configure")
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := testsPass(cwd)
	if err != nil {
		fmt.Println(err)
		fmt.Println(out)
		return
	}

	archive, err := createArchive(cwd)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created " + archive)

	// TODO
	err = uploadFile(config.ApiKey, archive)
	if err != nil {
		fmt.Println("Upload failed - ", err)
		return
	}
	fmt.Println("Successfully uploaded")
}

func testsPass(testDir string) (string, error) {
	cmd := exec.Command("go", "test")
	cmd.Dir = testDir
	out, err := cmd.Output()
	return string(out), err
}

func uploadFile(apikey, archive string) error {
	return nil
}
