This is the file containing our decision on how to contribute to our minitwit project

## Which repository setup will we use?
We choose the Mono repository setup.

## Which branching model will we use?
We have chosen a modification of gitflow, to make it a little less cumbersome we don't have a release branch.

We have a master containing code in production (it is treated like a release branch).
We have a development branch that is merged into the master to make releases.
We branch out from the development branch when we do a feature and then a PR is created, then a manager can merge the changes into the development branch. 

## Which distributed development workflow will we use?
We have decided to go with the Integration Manager workflow, where all of us acts ad managers. The managers are responsible for looking at PR's/commenting/approving/merging them.

## How do we expect contributions to look like?
1. ISSUES: We want to use Githubs project feature, to make issues. 
2. BRANCHES: We use these issue numbers as names for our branches e.g. "Issue 1"
3. COMMITS: Then a commit message for that branch would be "Issue 1: This is why the cat shouldn't sit on my keyboard."

## Who is responsible for integrating/reviewing contributions?
Then we push our feature/bug to the remote and create a pull request with a short but elaborate description of the content of the PR. One/two team members are added as reviewers, they have to leave comments and when they are ready they approve. When a PR has been approved anyone can merge it into the development branch. 