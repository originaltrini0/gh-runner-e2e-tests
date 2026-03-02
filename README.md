# GitHub Actions Runner E2E Tests

[![Self-hosted Runner pipelines](https://github.com/originaltrini0/gh-runner-e2e-tests/actions/workflows/e2e-tests.yml/badge.svg)](https://github.com/originaltrini0/gh-runner-e2e-tests/actions/workflows/e2e-tests.yml)

End-to-end tests for validating self-hosted GitHub Actions runners across multiple architectures.

## Overview

This repository contains a GitHub Actions workflow that exercises self-hosted runners to verify they can correctly handle common CI/CD workloads on both **amd64** and **arm64** architectures.

## Workflow: `e2e-tests.yml`

**Triggers:** `push` / `pull_request` to `main`/`master`, or manual `workflow_dispatch`

All jobs use a **matrix strategy** across `amd64` and `arm64`, running on self-hosted runners with the corresponding architecture label.

### Jobs

| Job | Stage | Description |
|-----|-------|-------------|
| **shell-tests** | — | Runs basic shell commands (`uname --all`) to confirm the runner is functional. (The Docker host's architecture, not the runner's architecture, is displayed. I need to test some more.)|
| **build-container** | — | Checks out the repo, builds a Docker image from the included `Dockerfile`, and runs the resulting container for the target platform. |
| **build-binary** | — | Checks out the repo, sets up Go 1.22, compiles `main.go` into a native binary for the target architecture, and executes it. |

### Runner Labels Required

Each runner must be registered with the labels:
- `self-hosted`
- `amd64` or `arm64` (matching its architecture)

## Supporting Files

| File | Purpose |
|------|---------|
| `Dockerfile` | Minimal Alpine image that writes and outputs the build architecture. Used by the **build-container** job. |
| `main.go` | Go program that reports its architecture and performs a simple computational workload (10M random float iterations). Used by the **build-binary** job. |
