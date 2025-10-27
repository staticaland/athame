# Renovate Module

A Dagger module for running Renovate to check for dependency updates in repositories.

## Usage

Run Renovate against a repository:

```bash
dagger call -m ./renovate run \
  --project="owner/repo" \
  --token=cmd://"gh auth token | tr -d '\n'"
```

**Important:** The `tr -d '\n'` is required to strip the trailing newline from the token, which would otherwise cause authentication failures.

### Parameters

- `project` - Repository to run Renovate against (e.g., `"staticaland/athame"`)
- `token` - GitHub authentication token (use `cmd://` to load from command)
- `platform` - Platform to use (defaults to `"github"`)

## Token Requirements

Your GitHub token needs the following scopes:

- `repo` - Full control of private repositories
- `workflow` - Update GitHub Action workflows (optional, for workflow updates)

Verify your token scopes with:

```bash
gh auth status
```

## Examples

### Using gh CLI token

```bash
dagger call -m ./renovate run \
  --project="staticaland/athame" \
  --token=cmd://"gh auth token | tr -d '\n'"
```

### Using environment variable

```bash
dagger call -m ./renovate run \
  --project="staticaland/athame" \
  --token=env:GITHUB_TOKEN
```

### Different platform

```bash
dagger call -m ./renovate run \
  --project="group/repo" \
  --token=env:GITLAB_TOKEN \
  --platform="gitlab"
```

## Base Container

Access the base Renovate container for custom workflows:

```bash
dagger call -m ./renovate base \
  with-exec --args="renovate,--version" \
  stdout
```
