# Handler: Dependabot Go Module Differential Sync

This sub-skill handles dependency updates from Dependabot by translating git diffs into clean local native Go commands.

**Execution policy:** All `skipper run …` steps run autonomously. Never ask before running them.

## Routing: two dependabot modes

Inspect the `go.mod` diff **before** acting:

| Diff signal | Handler mode |
|-------------|--------------|
| Changes root `go` directive **or** ≥2 `k8s.io/*` modules **or** `k8s.io/kubernetes` | **Toolchain mode** → follow `toolchain-handler.md` (Phase 2) in full |
| Single third-party module bump only (grpc, ginkgo, runc, …) | **Package mode** → steps below |

## Package mode (single dependency)

**CRITICAL SAFETY GUARDRAIL:** do not use `replace` directive, unless it was already defined in the go.mod for the specific package.
**CRITICAL SAFETY GUARDRAIL (USER INTERVENTION):** If an authorization token is required for pulling any image or dependency, STOP immediately, log the issue, and ask the user to provide authorization/tokens. Otherwise, proceed completely autonomously.

1. **Analyze Diffs:** Inspect the changes made specifically to the `go.mod` file inside the targeted external commit hash.
2. **Execute Go Fetch:** For each added or updated package version found in that diff, execute:
```bash
   skipper run go get <package>@<version>
   skipper run go mod tidy
   skipper run go mod vendor
   git commit -am "<original_commit_message>"
```
## Toolchain mode (Go / K8s upgrade)

When dependabot touches the toolchain, **do not** stop at `go get` for one package. Delegate entirely to `toolchain-handler.md`:
* Apply all module bumps from the diff (and any paired Markus `go.mod` commits if needed for consistency).
* Update `Dockerfile.node-feature-discovery-build` base image to match the new Go version.
* `skipper run go mod tidy`, `skipper run go mod vendor`, `skipper run make build`.
* Run Phase 3 compatibility-debt restoration from `compatibility-debt.md`.
