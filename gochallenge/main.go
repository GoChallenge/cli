package main

import (
	"os"

	"github.com/GoChallenge/cli/cmd"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go Challenge Submission Tool"
	app.Version = "0.1 alpha"
	app.Usage = "A tool to help programmers participate in the monthly Go challenge"
	app.Authors = []cli.Author{
		{
			Name: "https://github.com/GoChallenge/cli/graphs/contributors",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "configure",
			ShortName: "c",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key, k",
					Usage: "API Key from golang-challenge.com",
				}},
			Action: cmd.Configure,
			Usage:  "Writes the config values to a json file",
		},
		{
			Name:      "submit",
			ShortName: "s",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type, t",
					Usage: "Specify if you're participating normally or for fun",
					Value: "normal",
				}},
			Action: cmd.Submit,
			Usage:  "Submit your solution to the latest challenge",
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Action:    cmd.Fetch,
			Usage:     "Fetch the latest challenge",
		},
	}

	app.Run(os.Args)
}
