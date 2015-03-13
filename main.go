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
		// {
		// 	Name:      "fetch",
		// 	ShortName: "f",
		// 	Action:    fetch,
		// 	Usage:     "Fetch the latest challenge",
		// },
		{
			Name:      "login",
			ShortName: "l",
			Action:    login,
			Usage:     "Save golang-challenge.com api credentials",
		},
		{
			Name:      "logout",
			ShortName: "o",
			Action:    logout,
			Usage:     "Clear golang-challenge.com api credentials",
		},
		{
			Name:      "submit",
			ShortName: "s",
			Action:    submit,
			Usage:     "Submit your solution to the latest challenge",
		},
	}

	app.Run(os.Args)
}
