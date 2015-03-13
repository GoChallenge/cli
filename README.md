Go Challenge CLI
===============
The requirement is to create a web site and associated tools to allow easy participation in the GoChallenge. 

The process for the first GoChallenge was:

 - The participant downloaded a zip file
 - The participant wrote code to make those test cases pass
 - The entire Go package (the code written by the participant and the original files) was zipped up and emailed to the organiser, who then manually evaluated each submisssion.

From the second challenge onwards, participants will have the option of using a CLI (inspired by Exercism.io's CLI)

So far it has this functionality:

 - `configure` - Stores apikey locally. This is used during submission
 - `fetch` - Can retrieve a repo once it has the repo name. Uses the `go get` command to do this.
 - `submit` - Checks if apikey is present. If it is, runs the tests in the directory specified. If tests pass, makes a zip archive of the files/folder the user provides. 

TODO

 - configure should verify the apikey pair with the server
 - fetch should get the current repo name from the server (hardcoded for now)
 - fetch should retrieve past challenges
 - submit should upload the zip file to the server