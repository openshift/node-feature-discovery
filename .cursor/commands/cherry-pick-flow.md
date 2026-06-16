---
description: "Cherry-pick a sequential timeline of commits from an external repository into the current repository with iterative make validations."
---
# /cherry-pick

Use this command to safely pull a range of commits chronologically from an external repository into this workspace.

## Usage
`/cherry-pick <path_to_other_repo> <starting_commit_sha>`

## System Mapping
1. Extract the local folder path to the external repository from the first argument.
2. Extract the starting commit hash from the second argument.
3. Open, read, and strictly follow the execution instructions located inside `.cursor/skills/cross-repo-picker.md`.
4. Execute the loop **yourself**, one commit at a time — do **not** use helper scripts (`run-cherry-pick-queue.sh`, `process-cherry-pick.sh`, or similar batch automation).
5. Dispatch the execution loop passing the arguments:
   - Source Repo: $ARGUMENT_1
   - Start Commit: $ARGUMENT_2
