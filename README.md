# GitHub Actions Runner E2E Tests

[![Self-hosted Runner pipelines](https://github.com/originaltrini0/gh-runner-e2e-tests/actions/workflows/e2e-tests.yml/badge.svg)](https://github.com/originaltrini0/gh-runner-e2e-tests/actions/workflows/e2e-tests.yml)

End-to-end tests for validating an ARC (GitHub Actions Runner Controller) ephemeral runner scale set. Runner pods run in the `ci-builders` Kubernetes namespace and target the `github-runner-scale-set` scale set; Docker images are built with a shared remote BuildKit.

## Workflow: `e2e-tests.yml`

**Triggers:** `push` / `pull_request` to `main`/`master`, or manual `workflow_dispatch`

### Jobs

| Job | Description |
|-----|-------------|
| **shell-tests** | Prints the runner name, pod architecture (`uname -m` and `uname --all`), and the `docker buildx` client version. |
| **build-container** | Matrix over target platforms `amd64` and `arm64` (platform matrix, not runner selection). Checks out the repo, creates a buildx remote builder named `buildkit-github` pointing at the shared `buildkitd-github` BuildKit using client certs mounted at `/etc/buildkit`, then builds the included `Dockerfile` for `linux/<arch>` and tags the result `e2e-test-image-<arch>`. Because both matrix legs run concurrently this also exercises ARC scaling to 2 pods. |
| **build-binary** | Checks out the repo, sets up Go 1.22, builds `main.go` for the pod's native architecture and runs it, then cross-compiles the other architecture and verifies both binaries with `go version <binary>` (the cross-architecture binary is not executed). |

### Runner targeting

Jobs target ARC ephemeral runners with:

```yaml
runs-on: [github-runner-scale-set]
```

The scale-set name doubles as the runner label. Runner pods run in the `ci-builders` Kubernetes namespace; there are no architecture-labeled runners — all cluster nodes are `amd64`, and both `linux/amd64` and `linux/arm64` platforms are produced by the shared BuildKit (arm64 `RUN` steps execute under the node's QEMU emulation).

## Build backend: remote BuildKit (no Docker daemon)

Runner pods have no Docker daemon. Images are built with the buildx "remote" driver against a shared BuildKit instance:

- `BUILDKIT_HOST` — `tcp://buildkitd-github.ci-builders.svc.cluster.local:1234` (gRPC + TLS)
- Client certificates are mounted at `/etc/buildkit` by the ARC pod template: `ca.pem`, `cert.pem`, `key.pem`
- The buildx builder is created at runtime in the `build-container` job with `--driver remote` and the cert driver-opts

Because the builder is remote, `docker run` and `--load` are not available. Image behavior is instead verified by a `RUN` step inside the `Dockerfile` (`cat /test.txt`), which executes at build time on the target architecture.

## Supporting Files

| File | Purpose |
|------|---------|
| `Dockerfile` | Minimal Alpine image that writes and outputs the build architecture via `RUN`. Used by the **build-container** job. |
| `main.go` | Go program that reports its architecture and performs a simple computational workload (10M random float iterations). Used by the **build-binary** job. |
