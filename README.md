Go Challenge CLI
===============

The process for the first GoChallenge was:

 - The participant downloaded a zip file
 - The participant wrote code to make those test cases pass
 - The entire folder was zipped up and emailed to the organiser, who then manually evaluated each submisssion.

From the second challenge onwards, participants will have the option of using a CLI (inspired by exercism.io's CLI)

### To install:

 - Set up Go - http://golang.org/doc/install
 - Run `go get github.com/GoChallenge/cli` from the command line
 - `cd $GOPATH/src/GoChallenge/cli/gochallenge`
 - `go install`

### To use:

 - `gochallenge configure -k "key"` - Stores your key locally. Used when verifying your submission. Get your key from http://golang-challenge.com
 - `gochallenge list` - Lists all available challenges, both open and closed.
 - `gochallenge fetch` - Fetches the latest challenge. `-all` fetches all challenges, including older ones. `-challenge id` fetches the challenge with that id.
 - `gochallenge submit` - Runs the tests in the folder. If the tests pass, it uploads a .zip archive of the folder to the GoChallenge website.


### How to contribute

* [golang-challenge Slack channel](https://gophers.slack.com/messages/golang-challenge/) is the main
  discussion channel for this project - feel free to jump in and talk to people there.

* [Github issues](https://github.com/GoChallenge/cli/issues) -
  alternatively, if you have an idea that you want to discuss, feel free to open an
  issue on this project, and start the discussion.
