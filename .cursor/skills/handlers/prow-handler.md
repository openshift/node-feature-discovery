# Handler: Prow Robot Commit Skipper

This sub-skill handles automated CI/CD commits that should not be merged into the local repository.

## Execution Steps
1. Log a message in the chat informing the user that this commit belongs to the Prow Robot and is being skipped.
2. Do not execute any `git cherry-pick` or file modifications.
3. Advance the queue pointers immediately to the post-commit validation step in the main orchestrator.
