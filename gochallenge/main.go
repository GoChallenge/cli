package main

import (
	"fmt"
	"os"

	"github.com/GoChallenge/cli/cmd"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Go Challenge Submission Tool"
	app.Version = "0.1 beta"
	app.Usage = "A tool to help programmers participate in the monthly Go challenge"
	app.Authors = []cli.Author{
		{
			Name: "https://github.com/GoChallenge/cli/graphs/contributors",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server, s",
			Usage:  "The server host:port",
			EnvVar: "GOCHALLENGE_SERVER",
		},
	}
	app.Before = func(c *cli.Context) error {
		server := "gc.falsum.me"
		if c.IsSet("server") {
			server = c.String("server")
		}
		cmd.API_URL = fmt.Sprintf("http://%s/v1/challenges", server)
		return nil
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
				},
			},
			Action: cmd.Submit,
			Usage:  "Submit your solution to the latest challenge",
		},
		{
			Name:      "fetch",
			ShortName: "f",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "challenge, c",
					Usage: "Fetch specified challenge",
				},
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Fetch all challenges -- Use carefully!",
				},
			},
			Action: cmd.Fetch,
			Usage:  "Fetch the latest challenge",
		},
		{
			Name:      "list",
			ShortName: "l",
			Action:    cmd.List,
			Usage:     "List all available challenges",
		},
	}

	app.Run(os.Args)
}
