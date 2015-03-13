Go Challenge CLI
===============
The requirement is to create a web site and associated tools to allow easy participation in the GoChallenge. 

The process for the first GoChallenge was:

 - The participant downloaded a zip file
 - The participant wrote code to make those test cases pass
 - The entire Go package (the code written by the participant and the original files) was zipped up and emailed to the organiser, who then manually evaluated each submisssion.

From the second challenge onwards, participants will have the option of using a CLI (inspired by Exercism.io's CLI)

So far it has this functionality:

 - `login` - Participants can store their username and apikey locally. This info is used during submission
 - `logout` - Stored info is deleted
 - `fetch` - Can retrieve a repo once it has the repo name. Uses the `go get` command to do this.
 - `submit` - Checks if the user is logged in. If he is, makes a zip archive of the files/folder the user provides. 

TODO

 - login should verify the username+apikey pair with the server
 - fetch should get the current repo name from the server (hardcoded for now)
 - fetch should retrieve past challenges
 - submit should upload the zip file to the server
 - submit should run the tests