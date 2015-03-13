package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
)

func fetch(c *cli.Context) {
	// TODO get the repo name from the server

	reponame := "github.com/GoChallenge/challenges/"

	err := exec.Command("go", "get", reponame).Run()
	if err != nil {
		// TODO this returns 1 even if it succeeds
		fmt.Println("Could not fetch the challenge", err)
		return
	}
	repopath := path.Join(os.Getenv("GOPATH"), "src", reponame)
	fmt.Println("Latest challenge downloaded. You can find it at", repopath)
}
