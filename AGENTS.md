# Node Feature Discovery (OpenShift fork)

Kubernetes add-on that detects hardware features and system configuration on nodes and advertises them as node labels, annotations, and extended resources. This is the OpenShift downstream fork of [kubernetes-sigs/node-feature-discovery](https://github.com/kubernetes-sigs/node-feature-discovery).

## Dev Environment

- Go 1.25+ (see `go.mod`)
- `podman` (or `docker` via `CONTAINER_COMMAND=docker`)
- `kubectl`, `helm` for deployment workflows
- `golangci-lint` for linting (installed by CI via `scripts/test-infra/verify.sh`)

## Build

```sh
make build          # compile all binaries into bin/
make image          # build container image (uses podman by default)
make install        # go install all binaries
```

Binaries are built from `cmd/`:
- `nfd-master` — cluster-side component, manages node labels
- `nfd-worker` — node-side component, detects features
- `nfd-topology-updater` — reports NUMA topology
- `nfd-gc` — garbage collector for stale NFD labels
- `nfd` — unified binary with subcommands
- `kubectl-nfd` — kubectl plugin

Override variables: `IMAGE_REGISTRY`, `IMAGE_TAG_NAME`, `CONTAINER_COMMAND`, `HOSTMOUNT_PREFIX`.

## Test

Run unit tests:

```sh
make test
```

This runs `go test -covermode=atomic -coverprofile=coverage.out ./cmd/... ./pkg/... ./source/...` and separately tests `api/nfd/...`.

Test files live alongside source files (e.g. `pkg/nfd-master/nfd-master_test.go`).

Run e2e tests (requires a running cluster and `KUBECONFIG`):

```sh
make e2e-test KUBECONFIG=<path>
```

E2e tests are in `test/e2e/` and use Ginkgo v2 (`github.com/onsi/ginkgo/v2`). Test data is in `test/e2e/data/`.

## Code Style

- Format: `gofmt -s` — run `make verify-gofmt` to check
- Lint: `golangci-lint run --timeout 10m` via `make ci-lint`
- Markdown lint: `make mdlint` (runs in a container)
- Helm lint: `make helm-lint`
- Log checker: `logcheck` with config at `scripts/test-infra/logcheck.conf`
- Boilerplate header: `hack/boilerplate.go.txt`

CI runs all checks via `scripts/test-infra/verify.sh`: gofmt, golangci-lint, helm lint, logcheck, unit tests, template freshness, kustomize overlay validation, Helm schema and README sync.

## Architecture

### Project Layout

- `cmd/` — entry points for each binary
- `pkg/` — core packages:
  - `pkg/nfd-master/` — master server logic, label management, API controller
  - `pkg/nfd-worker/` — worker logic, feature detection orchestration
  - `pkg/nfd-topology-updater/` — topology updater logic
  - `pkg/nfd-gc/` — garbage collector
  - `api/nfd/v1alpha1/types.go` — CRD type definitions (`NodeFeatureRule`)
  - `pkg/features/` — feature representation types
  - `pkg/resourcemonitor/` — resource monitoring and pod resource scanning
  - `pkg/utils/` — shared utilities (flags, hostpath, kubernetes helpers)
- `source/` — feature source plugins (cpu, kernel, memory, network, pci, storage, system, usb, local, custom, fake)
- `api/nfd/` — separate Go module for the NFD API types
- `deployment/` — Kubernetes manifests:
  - `deployment/base/` — kustomize base resources
  - `deployment/overlays/` — kustomize overlays (default, prune, prometheus, etc.)
  - `deployment/helm/node-feature-discovery/` — Helm chart
  - `deployment/components/` — component configs (worker, master, topology-updater)
- `hack/` — code generation and release scripts
- `scripts/test-infra/` — CI scripts (build, verify, e2e)

### Key Abstractions

Feature sources in `source/` implement the `LabelSource` interface (see `source/source.go`). Each source detects specific hardware/OS features (CPU flags, kernel config, PCI devices, etc.) and returns them as labels.

The master uses a `NodeFeatureRule` CRD (defined in `pkg/apis/nfd/`) to apply custom labeling rules.

Communication between worker and master using NodeFeature object.

### Multi-module Structure

The project has two Go modules:
- Root module: `github.com/openshift/node-feature-discovery` (Go 1.25)
- API module: `github.com/openshift/node-feature-discovery/api/nfd` (Go 1.24, separate `go.mod`)

The root module replaces the API module with a local path: `replace github.com/openshift/node-feature-discovery/api/nfd => ./api/nfd`.

## Deployment

Generate kustomize YAMLs:

```sh
make yamls
```

Deploy to a cluster:

```sh
make deploy
```

Clean NFD labels:

```sh
make deploy-prune
```

Update Helm values from config examples:

```sh
make templates
```

Generate CRDs and client code:

```sh
make generate
```
