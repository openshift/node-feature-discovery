# Handler: Standard Cherry-Picker & Conflict Resolver

This sub-skill executes traditional cherry-picks and guides the agent through conflict mitigation.

**Execution policy:** Cherry-pick and post-resolution `skipper run make build` run autonomously. Never ask before running skipper.

## Execution Steps

### 1. Pre-Pick File Evaluation
Before or during the cherry-pick process, you must strictly enforce the following repository-specific rules:
* **Rule 1 (Ignore Makefile Changes):** Any modifications, additions, or deletions to the `Makefile` must be completely ignored. Do not bring `Makefile` changes into this workspace.
* **Rule 2 (Ignore Non-Source Files):** Completely ignore any incoming files that are not core source files. This includes, but is not limited to:
  * Documentation files (e.g., `.md`, `.txt`, `/docs/` folders)
  * Test/Mock generation scripts (e.g., scripts for `mockgen`, vector generation, or test automation helpers)
* **Rule 3 (Phase 1 — compatibility, not toolchain upgrade):** If the incoming patch uses an API, package, or generated code that requires a **newer Go or K8s version** than the workspace currently has:
  1. **Do not** bump `go` in `go.mod`, `k8s.io/*` versions, or `Dockerfile.node-feature-discovery-build` in this handler.
  2. **Do** substitute an equivalent that compiles on the current toolchain (older function, revert generated-code hunk to pre-bump style, drop a hunk that only exists for newer client-go, etc.).
  3. **Record** each substitution in `.cursor/compatibility-debt.jsonl` per `compatibility-debt.md`.
  4. Log: `COMPAT-DEBT: <file> — deferred until toolchain upgrade`.
  5. The toolchain-handler will land the real version bump later; compatibility-debt restoration (Phase 3) will then switch back to the upstream API.

### 2. Execute Pick
* Run `git cherry-pick <SHA>` to attempt bringing the patch into this workspace.

### 3. Rewrite upstream import paths
After the cherry-pick lands (cleanly or after conflict resolution), scan all files touched by the commit for Go imports starting with `sigs.k8s.io/node-feature-discovery/`. Replace them with `github.com/openshift/node-feature-discovery/`. This applies to both production and test files.

Quick check:
```bash
git diff HEAD~1 --name-only -- '*.go' | xargs grep -l 'sigs.k8s.io/node-feature-discovery/' 2>/dev/null
```
For each match, replace `sigs.k8s.io/node-feature-discovery/` with `github.com/openshift/node-feature-discovery/` in the import block, stage the file, and amend the cherry-pick commit.

### 4. Handle Potential Merge Conflicts
* If the cherry-pick succeeds cleanly, exit this handler successfully.
* If the cherry-pick halts due to a **Merge Conflict**:
  1. **Log the Conflict:** Print details of the flagged files and the nature of the conflict in the chat log.
  2. **Attempt Resolution:** Intelligently analyze and attempt to resolve the source conflicts matching the contextual logic of the incoming patch. Apply Rule 3 when resolution would require a toolchain bump.
  3. **Conflict resolution (autonomous):** When multiple resolutions exist, choose the one that preserves upstream patch intent, downstream import paths, and handler ignore-rules (Makefile, docs, helm). Do **not** ask the user to pick.
  4. Stage resolutions: `git add <resolved_files>` and run `git cherry-pick --continue`.
  5. Run `skipper run make build` without asking.
  6. **Ask the user only when**:
     * **Auth required:** registry/token/credential needed for image or module pull.
     * **Failure after retry:** `cherry-pick --continue` or `skipper run make build` fails twice for this commit despite a reasoned fix attempt.
