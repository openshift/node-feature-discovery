# Handler: Go / K8s Toolchain Upgrade

Handles commits whose **primary purpose** is bumping the Go toolchain or coordinated Kubernetes module versions — including:
* `Bump Golang to vX.Y` (e.g. Markus Lehtonen, `3c5951cf`)
* `go.mod: update all k8s modules to vX.Y.Z`
* `go.mod: bump kubernetes to vX.Y.Z`
* Dependabot commits that change the root `go` directive **or** multiple `k8s.io/*` modules in one diff

## Phase 2 — Full toolchain upgrade (do everything needed)

**CRITICAL:** A toolchain upgrade is never just a `go.mod` one-liner. Complete **all** of the following before committing:

### 1. Apply version bumps
* Mirror commit's `go.mod` diff with `skipper run go get` for each changed module.
* Update `api/nfd/go.mod` if the upstream commit touches it.

### 2. Update build images (mandatory)
* Edit `Dockerfile.node-feature-discovery-build` so the `FROM` builder image matches the new Go version (e.g. `golang-1.26-openshift-…`).
* **Do not** touch production Dockerfile, it will be updated manually by the user
* **Do not** rely on upstream `Makefile` `BUILDER_IMAGE` defaults — downstream uses OpenShift builder tags; Makefile changes are ignored.

#### Missing builder image — STOP and ask
If the builder image for the new Go version does not exist in the registry (e.g. `manifest unknown`, 404, or image-pull error during `skipper run`):
1. **Do NOT** fall back to the previous builder image or revert the Go version.
2. **Do NOT** skip the commit.
3. **STOP** and ask the user for the correct builder image tag. The user knows which images and OpenShift version suffixes are available.
4. Resume the toolchain upgrade only after the user provides the correct image.

### 3. Vendor and verify
```bash
skipper run go mod tidy
skipper run go mod vendor
skipper run make build
```

### 4. Commit
```bash
git commit -am "<original_commit_message>"
```

## Failure handling
* If `skipper run go mod tidy` / `skipper run go get` pulls incompatible k8s module sets, apply the **coordinated** k8s bump commits together (dependabot k8sio group + Markus `go.mod` commits) in one toolchain session — do not land a partial k8s upgrade.
* If build fails after image update, fix autonomously and re-run `skipper run make build` once. Ask the user only if the retry fails or the error indicates missing registry/auth credentials.

## What to ignore from upstream toolchain commits
* `Makefile` / `Tiltfile` builder-image defaults (Rule 1 in standard-picker) — already covered by `Dockerfile.node-feature-discovery-build` update above.
