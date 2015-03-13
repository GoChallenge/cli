Go Challenge CLI
===============
The requirement is to create a web site and associated tools to allow easy participation in the GoChallenge. This issue is to discuss a CLI tool that can be used by the participants to interact with the site.

Right now, the process for the GoChallenge is:
 - The problem is posted to the site
 - The participant downloads a file or clones a repo that contains a few test cases
 - The participant writes code to make those test cases pass
 - The entire Go package (the code written by the participant and the original files) is zipped up and emailed to the organiser, who then manually checks the submission and marks it as successful or failed

Taking inspiration from the CLI tool used by the Exercism.io site, the GoChallenge CLI should at least:
 - Allow the participant to connect their installation of the tool to their account on the GoChallenge web site
 - Download the files for the currently running challenge task
 - Once the participant has written code that can pass the test cases, the tool should allow them to submit it to the GoChallenge website; after making sure the test cases are passing. While not needed, ensuring that the test cases pass before submitting the code avoids accidental submissions

Once the code has been submitted, there should appear on the users account info on the GoChallenge site a submission entry with links to download the associated files. While the challenge is still on-going, only the challenge organiser or the participant should be able to download that file. Once the challenge has been completed, the file should be visible to everyone, so others can look and learn from the code.

We might also want to include functionality on the web site to automatically run the submitted code and verify that the test cases do indeed pass. In this manner, the organiser will not need to go through the submissions manually and can verify their correctness later on, towards the end of the challenge. But given the security implications of running untrusted code automatically, we might want to hold off on this. For now, I think an email to the organiser when a participant submits their code should be fine.

Please add any required feature that I might have missed, and lets start a discussion on this so we can have a tool that will be really useful to others.
