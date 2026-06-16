# Skill: Iterative Cross-Repository Cherry-Picker

This skill executes a step-by-step sequential cherry-pick loop from an external repository while maintaining local workspace safety.

## Three-phase compatibility model

Cherry-picking across different Go/K8s baselines uses three phases. Handlers must respect phase boundaries:

| Phase | When | Handler | Action |
|-------|------|---------|--------|
| **1 — Defer** | Standard cherry-pick needs newer Go/K8s API | `standard-picker-handler.md` | Use older compatible API; record in `.cursor/compatibility-debt.jsonl` |
| **2 — Upgrade** | Commit explicitly bumps Go/K8s (toolchain or dependabot k8sio) | `toolchain-handler.md` / dependabot toolchain mode | Full upgrade incl. `Dockerfile.node-feature-discovery-build` |
| **3 — Restore** | Immediately after Phase 2 build passes | `compatibility-debt.md` | Re-apply upstream APIs from debt ledger; build & commit |

## Global Execution Rules
* **One commit at a time — no helper scripts:** Process the queue **manually**, commit by commit, by reading the relevant handler skill and running git/skipper commands yourself.
* **Never** invoke batch scripts (e.g. `run-cherry-pick-queue.sh`, `process-cherry-pick.sh`) or write new automation to loop through the queue. After each commit is applied or skipped, run validation, then move to the **next single commit** — do not batch cherry-picks or builds across multiple commits.
* **Verbose Logging:** Log every atomic action, decision, skipped commit, or successful merge in the chat console in real-time. Do not execute silently.
* **Debt logging:** Any Phase 1 downgrade must log `COMPAT-DEBT: …`. Any Phase 3 restoration must log `COMPAT-RESTORE: …`.

## Step 1: Commit Range Extraction
1. Navigate to the external repository path provided by the user.
2.  Generate an ordered list of all commits from `<starting_commit_sha>` up to `HEAD`, in the exact order they appear when scrolling **up** in `git log` output.
   * *Formula:* Run `git log --format="%H|%an|%s"` (no other flags) in the source repo, pipe the output to a file, find the line containing `<starting_commit_sha>`, take that line and all lines above it (toward HEAD), then reverse them so `<starting_commit_sha>` is first and `HEAD` is last. Shell one-liner: `git log --format="%H|%an|%s" | head -n $(git log --format="%H" | grep -n <starting_commit_sha> | cut -d: -f1) | tac`
3. Parse this list into a clean state array containing `[SHA, Author, Commit Message]`. 
4. Present the list to the user and wait for his approval to proceed.

## Step 2: Prepare Target Directory
1. Navigate to target directory
2. Check if `sig` target is defined in the .git/config. If it is not defined, add it to the .git/config using `git remote add sig https://github.com/kubernetes-sigs/node-feature-discovery`
3. Run `git fetch sig`
4. Ensure `.cursor/compatibility-debt.jsonl` exists (create empty if missing).

## Step 3: The Iterative Execution Loop
Process **exactly one commit** per iteration. For that commit, evaluate the **Author** and **commit message** and route:

### Use Case 1: Prow Robot (Skip)
* **Condition:** Author is `"Kubernetes Prow Robot"`.
* **Action:** Delegate immediately to `.cursor/skills/handlers/prow-handler.md`.

### Use Case 2: Toolchain upgrade (Phase 2)
* **Condition:** Commit message matches `Bump Golang to`, `go.mod: update all k8s`, `go.mod: bump kubernetes`, **or** dependabot diff changes root `go` / multiple `k8s.io/*` (see `dependabot-handler.md` routing).
* **Action:** Delegate to `.cursor/skills/handlers/toolchain-handler.md`, then run Phase 3 from `compatibility-debt.md`.

### Use Case 3: Dependabot (package or toolchain)
* **Condition:** Author is `"dependabot[bot]"`.
* **Action:** Delegate to `.cursor/skills/handlers/dependabot-handler.md` (routes to package mode or toolchain mode internally).

### Use Case 4: Standard Commits (Cherry-Pick & Resolve, Phase 1)
* **Condition:** Any other author.
* **Action:** Delegate to `.cursor/skills/handlers/standard-picker-handler.md`.

---

## Step 4: Post-Commit Validation & Pacing
Immediately after a handler successfully completes a commit lifecycle (or a skip occurs):

1. **Sanity Check:** If a code change occurred, check that code compiles:
```bash
   skipper run make build
```
   On failure: fix once and retry; ask the user only if the retry fails or auth is required.

2. **Phase 3 check:** If the commit just applied was a toolchain upgrade (Use Case 2 or dependabot toolchain mode), scan `.cursor/compatibility-debt.jsonl` and restore all pending entries before advancing the queue. Execute `Sanity Check` before advancing the queue.

3. **Advance:** Log the result (`APPLIED`, `SKIPPED`, or `FAILED`) for this commit, then return to Step 3 for the **next** commit only. Repeat until the queue is exhausted or a hard failure requires user input.

## Step 5: Final Test Run

After the queue is fully exhausted and the final `skipper run make build` has passed, run the full test suite:

```bash
skipper run make test
```

* If all tests pass, proceed to Step 6.
* If tests fail, analyze the failures, fix them (import path mismatches, compatibility issues, etc.), re-run `skipper run make test`, and commit any fixes with a descriptive message before proceeding.
* Ask the user only if tests still fail after one fix attempt.

## Step 6: Cleanup

Once the queue is fully exhausted (all commits applied or skipped), delete the runtime artifacts created during this run:

```bash
rm -f .cursor/compatibility-debt.jsonl
rm -f .cursor/cherry-pick-queue.txt
rm -f .cursor/cherry-pick-progress.txt
```

Do **not** delete the skill/handler markdown files or the command file — they are reusable for future runs.
