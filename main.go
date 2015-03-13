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
			Name:      "configure",
			ShortName: "c",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key, k",
					Usage: "exercism.io API key (see http://exercism.io/account)",
				}},
			Action: configure,
			Usage:  "Writes the config values to a json file",
		},
		{
			Name:      "submit",
			ShortName: "s",
			Action:    submit,
			Usage:     "Submit your solution to the latest challenge",
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Action:    fetch,
			Usage:     "Fetch the latest challenge",
		},
	}

	app.Run(os.Args)
}
