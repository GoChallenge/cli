package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
)

func List(c *cli.Context) {
	challenges, err := getChallengeDescriptors(API_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)
	for _, chal := range challenges {
		fmt.Fprintf(w, "%s\t(%s)\t   ---  \tid: %d\n", chal.Name, chal.Status, chal.ID)
	}
	w.Flush()
}
