Go Challenge CLI
===============
The requirement is to create a web site and associated tools to allow easy participation in the GoChallenge. 

The process for the first GoChallenge was:

 - The participant downloaded a zip file
 - The participant wrote code to make those test cases pass
 - The entire Go package (the code written by the participant and the original files) was zipped up and emailed to the organiser, who then manually evaluated each submisssion.

From the second challenge onwards, participants will have the option of using a CLI (inspired by exercism.io's CLI)

So far it has this functionality:

 - `configure` - Stores apikey locally. This is used during submission
 - `fetch` - Can retrieve a repo once it has the repo name. Uses the `go get` command to do this.
 - `submit` - Checks if apikey and challenge is present. If they are, runs the tests in the directory specified. If tests pass, makes a zip archive of the files/folder the user provides and uploads the file.

TODO

 - fetch should retrieve past challenges

Testing:
 This repo contains a fake_api.json file that can be used to test a few aspects of the CLI tool. To run it, follow these steps. You'll need nodejs and NPM installed:
 - npm install -g json-server
 - json-server fake_api.json
 - Enjoy!
