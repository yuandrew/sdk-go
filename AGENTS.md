# Contributor Guide

This repository contains the Temporal Go SDK. The following notes provide an overview of the repo layout, expectations for code changes, and how to run automated checks.

## Repo layout

- `client/`, `workflow/`, `activity/`, `temporal/`, and `worker/` hold the public API implementation.
- `internal/` contains private logic and utilities used by the public packages.
- `test/` contains integration and unit test files.
- `testsuite/` offers helpers such as a dev server for integration tests.
- `contrib/` holds optional utilities including the workflow determinism checker under `contrib/tools/workflowcheck`.
- `internal/cmd/build/` provides a helper program used to run static checks and tests.
- The contributor guide lives at `CONTRIBUTING.md`.

## Static checks and testing

The repo provides a build helper to run standard checks.
Change to `internal/cmd/build` and execute the following as needed:

```bash
# Run go vet, errcheck, staticcheck, and documentation link checks
go run . check

# Run integration tests (pass -dev-server to start a local Temporal server)
go run . integration-test

# Run unit tests
go run . unit-test
```

`integration-test` and `unit-test` accept `-run` to filter tests. Coverage output can be enabled with `-coverage-file` or `-coverage`.

## Utilities

- `workflowcheck` in `contrib/tools/workflowcheck` statically detects workflow non-determinism. Install with `go install` or run via the `workflowcheck` command.
- `doclink` in `internal/cmd/tools/doclink` ensures exported symbols that wrap internal ones contain proper documentation links. The `check` command runs this automatically.

## Pull request expectations

- Follow the commit message style described in `CONTRIBUTING.md`. PR titles become commit messages and must follow the [Chris Beams](http://chris.beams.io/posts/git-commit/) conventions. Avoid generic titles and start each title with an uppercase letter.
- Keep changes focused and provide a clear description of what and why in the PR body.
- Ensure all static analysis and tests pass by running the commands above before submitting a PR.
- PRs should not modify generated files.

## Reviewer checklist

Reviewers will verify that:

- `go run . check` passes with no issues.
- All relevant unit and integration tests pass.
- Commit messages and PR title follow the style guide.
- New code includes appropriate comments and uses existing packages where possible.
- The PR is focused and does not introduce unrelated changes.

Anything else you might not understand, refer to the README.md
