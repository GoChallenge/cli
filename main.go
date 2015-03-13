package main

import (
	"os"

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
			Name:      "config",
			ShortName: "c",
			Action:    writeConfig,
			Usage:     "Save and displayes configuration settings for the CLI app",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "apiKey",
					Usage: "Your api key. Get this from golang-challenge.com/account/. If this flag isn't present, prints out the current config",
				},
			},
		},
	}

	app.Run(os.Args)
}
