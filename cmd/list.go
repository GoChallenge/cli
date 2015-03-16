package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func List(c *cli.Context) {
	challenges, err := getChallengeDescriptors(apiURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(challenges) == 0 {
		fmt.Println("Found no challenges")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)
	for _, chal := range challenges {
		fmt.Fprintf(w, "%s\t(%s)\t   ---  \tid: %d\n", chal.Name, chal.Status, chal.ID)
	}
	w.Flush()
}
