# Contributing

## Development setup

Requires [Homebrew](https://brew.sh) (macOS) and [Go](https://go.dev).

```bash
make setup   # brew bundle + trunk install + git aliases
```

This installs all required tools (`go`, `terraform`, `golangci-lint`, `gh`, `trunk`) and wires up the repo-local git aliases in `.gitconfig`.

## Commit conventions

This repo uses [Conventional Commits](https://www.conventionalcommits.org). Commitlint enforces this via a Trunk pre-commit hook after `make setup`.

| Type                                     | Triggers release |
| ---------------------------------------- | ---------------- |
| `feat`                                   | minor            |
| `fix`, `perf`, `revert`, `docs`, `chore` | patch            |
| `build`                                  | minor            |
| `style`, `test`, `ci`                    | none             |
| `BREAKING CHANGE` footer                 | major            |

## Branching strategy

| Branch       | Purpose                                     |
| ------------ | ------------------------------------------- |
| `main`       | Stable releases — merge here for production |
| `beta`       | Beta pre-releases (`v1.0.0-beta.1`)         |
| `alpha`      | Alpha pre-releases (`v1.0.0-alpha.1`)       |
| `next`       | Next minor pre-releases (`v1.1.0-next.1`)   |
| `next-major` | Next major pre-releases (`v2.0.0-next.1`)   |

## Releases

Releases are fully automated via [semantic-release](https://semantic-release.gitbook.io) and [GoReleaser](https://goreleaser.com).

### Stable release

Merge a PR into `main`. Semantic-release analyzes commit messages since the last tag and creates a new version tag if warranted. GoReleaser builds and publishes to GitHub Releases and the [Terraform Registry](https://registry.terraform.io/providers/thomasgeens/resourcenamingtool).

### Pre-release

1. Create a pre-release branch (e.g. `beta`) off `main` if it doesn't exist.
2. Merge feature PRs into that branch instead of `main`.
3. Semantic-release creates pre-release tags automatically (e.g. `v1.0.0-beta.1`).
4. GoReleaser publishes a GitHub pre-release.

Pre-release versions are available on the Terraform Registry but are not selected by open-ended version constraints (`~> 1.0`). Users must pin explicitly:

```hcl
terraform {
  required_providers {
    resourcenamingtool = {
      source  = "thomasgeens/resourcenamingtool"
      version = "1.0.0-beta.1"
    }
  }
}
```

When the pre-release is stable, merge the pre-release branch into `main` to cut a stable release.

## Running tests

```bash
make test      # unit tests (TF_ACC=0)
make testacc   # acceptance tests (TF_ACC=1, requires Terraform)
make lint      # golangci-lint
make generate  # regenerate docs — commit the result
```
