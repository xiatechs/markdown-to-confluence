# **How to contribute**

---
The following is a set of guidelines for contributing to our tooling projects. These are mostly guidlelines, not rules, so use your best judgement and feel free to propose changes to this document in a pull request.

<br />  

### **First steps**
The first thing to do is check the [jira board](https://xiatech.atlassian.net/jira/software/projects/XT/boards/106) (filtered by the relevant label) to see what work has already been planned for the project. If you're not sure what to tackle first, reach out to the [xiatech-tooling slack channel](https://slack.com/app_redirect?channel=xiatech-tooling) and they should be able to point you in the right direction.

<br />

### **Want to add a new feature?**
The first thing to do is check to see if it is on the jira board. If it is and hasn't been started then assign the ticket to yourself, move it to `In Progress`, create a new code branch and away you go.

There may be cases where your feature is already on the board and has already been picked up by someone else. If so, reach out to them on the Xiatech-Tooling slack channel.

If your feature doesn't exist on the board, create a ticket, adding as much detail as possible, assign it to yourself and get started.

<br />  

### **Did you find a bug?**
Check the jira board to see if it has already been reported. If it is being worked on, reach out to whoever is working on it via the Xiatech-Tooling slack channel.

If it isn't on the board, create a defect with as much information as possible. If you're not able to work on it for any reason, post a message in the Xiatech-Tooling slack channel to let others know that you've found an issue.

<br />  

### **Submitting**
Always write clear messages for your commits. One line messages are fine for small, simple changes but more information is required for larger, more complex changes.

When creating a Pull Request, thoroughly fill out the PR template. Once completed, post a link to your PR in the Xiatech-Tooling slack channel.

<br />  

### **Code Conventions**
* Ensure your code passes a linter. We use `golangci-lint` and only merge in to the main branch if there are no linter errors

* Tests. We like unit tests. Add them, and be thorough with your test cases

* We try to leave the codebase in a better state than we found it so if you see that something needs refactoring, make the change.