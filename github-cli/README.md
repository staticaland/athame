# GitHub CLI Module

A Dagger module for working with GitHub CLI (gh) via asdf version manager.

## Usage

### List Repositories

List repositories for the authenticated user:

```bash
dagger call -m ./github-cli list-repos \
  --token=cmd://"gh auth token | tr -d '\n'"
```

**Important:** The `tr -d '\n'` is required to strip the trailing newline from the token, which would otherwise cause authentication failures.

### Parameters

- `token` - GitHub authentication token (use `cmd://` to load from command)
- `limit` - Maximum number of repositories to list (defaults to `"100"`)

## Base Container

Access the base GitHub CLI container for custom workflows:

```bash
dagger call -m ./github-cli base \
  with-exec --args="gh,--version" \
  stdout
```
