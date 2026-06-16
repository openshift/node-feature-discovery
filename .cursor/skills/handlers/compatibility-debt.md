# Ledger: Compatibility Debt (API Downgrades)

Tracks source-code adaptations made during cherry-picks **before** the matching Go/K8s toolchain upgrade lands. Each entry is reverted in Phase 3 once the toolchain commit has been applied and `skipper run make build` passes.

## Ledger file
Maintain at: `.cursor/compatibility-debt.jsonl`

Each line is one JSON object:
```json
{"commit":"<upstream-sha-being-picked>","file":"path/to/file.go","reason":"API X requires go 1.26","upstream_api":"cache.ToListWatcherWithWatchListSemantics","local_api":"&cache.ListWatch{","restore_after":"go-toolchain-upgrade"}
```

## Phase 1 — Record debt (standard-picker-handler)
When you downgrade or substitute an API/package call so the patch compiles on the **current** Go/K8s toolchain:
1. Apply the compatible alternative.
2. Append one ledger entry per adaptation.
3. Log: `COMPAT-DEBT: recorded <file> — will restore after toolchain upgrade`.

## Phase 3 — Restore debt (post toolchain upgrade)
After any handler completes a Go or K8s toolchain upgrade and build passes:
1. Read `.cursor/compatibility-debt.jsonl`.
2. For each open entry, re-apply the **upstream** API/version from the original patch intent.
3. Run `skipper run make build` autonomously.
4. On success, remove the entry from the ledger. On failure after one fix attempt, ask the user.
5. Commit restorations separately or amend the toolchain commit — prefer a dedicated commit: `restore: re-enable deferred APIs after toolchain upgrade`.

## When NOT to record debt
* Import-path rewrites (`sigs.k8s.io` → `github.com/openshift/...`) — permanent downstream difference.
* Intentionally ignored files (Makefile, docs, helm) — not debt.
* Prow-skip or helm-skip — no code landed.
