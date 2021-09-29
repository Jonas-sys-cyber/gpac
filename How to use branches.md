 * Create a Branch: `git branch branchname`
 * Switch to Branch: `git switch branchname`
 * Show current branch and all available: `git branch`
 * Delete Branch: `git branch -d branchname`
 * Merge a dev Branch with master:
    * First, switch to master: `git switch master`
    * Then, merge it: `git merge dev`
    * Look with `git status` if you have a clean working dir
    * If it is fine, delete the dev branch if you are finished: `git branch -d dev`
 