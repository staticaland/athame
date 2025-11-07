# MkDocs CI

A Dagger module for MkDocs CI/CD: linting, building, and publishing documentation sites.

This module provides functions to lint markdown with vale, prettier, and markdownlint-cli2, build MkDocs Material sites, and publish the site as a container image to GitHub Container Registry (GHCR).

## Features

- **Linting**: Run vale, prettier, and markdownlint-cli2 on your markdown files
- **Link checking**: Validate links with lychee
- **Concurrent testing**: All linters run in parallel for fast feedback
- **Build**: Build MkDocs Material sites
- **Publish**: Publish sites as container images to GitHub Container Registry (GHCR)
- **Deploy**: Deploy to Fly.io and/or trigger Render deploy hooks after successful publish
- **Notifications**: Send ntfy notifications at key pipeline stages (start, tests done, deployment complete)

## Usage

### Run individual linters

```bash
# Run vale
dagger call --mod ./mkdocs-ci vale

# Run prettier
dagger call --mod ./mkdocs-ci prettier

# Run markdownlint
dagger call --mod ./mkdocs-ci markdownlint

# Run link checking
dagger call --mod ./mkdocs-ci check-links
```

### Run all tests concurrently

```bash
dagger call --mod ./mkdocs-ci run-all-tests
```

### Build the site

```bash
dagger call --mod ./mkdocs-ci build export --path=./output
```

### Publish to GitHub Container Registry

To publish to GHCR, you need a GitHub Personal Access Token with `write:packages` permission. Use the GitHub CLI to get your token:

```bash
dagger call --mod ./mkdocs-ci publish \
  --ghcr-token=cmd:"gh auth token | tr -d '\n'"
```

By default, this publishes to `ghcr.io/staticaland/athame/<image-name>:latest`. You can customize:

```bash
dagger call --mod ./mkdocs-ci publish \
  --image-name=my-docs \
  --tag=v1.0.0 \
  --ghcr-username=myusername \
  --ghcr-token=cmd:"gh auth token | tr -d '\n'"
```

The publish function will:

1. Build the MkDocs Material site
2. Create a multi-platform container image (linux/amd64 and linux/arm64) with nginx and the static files
3. Publish to GitHub Container Registry

### Deploy to Fly.io

```bash
dagger call --mod ./mkdocs-ci deploy-flyio \
  --app=my-app-name \
  --image="ghcr.io/staticaland/athame/mkdocs-demo:latest@sha256:..." \
  --token=env:FLY_API_TOKEN
```

This deploys a published container image to Fly.io. The Fly.io API token should be stored as an environment variable.

### Trigger Render deploy

```bash
dagger call --mod ./mkdocs-ci deploy-render \
  --deploy-hook-url=env:RENDER_DEPLOY_HOOK_URL
```

This triggers a Render deploy hook. The deploy hook URL should be stored as a secret (e.g., environment variable).

### Full CI/CD pipeline

Run all linters, build, and publish if tests pass:

```bash
dagger call --mod ./mkdocs-ci lint-build-publish \
  --ghcr-token=cmd:"gh auth token | tr -d '\n'"
```

With Render deploy hook integration:

```bash
dagger call --mod ./mkdocs-ci lint-build-publish \
  --ghcr-token=cmd:"gh auth token | tr -d '\n'" \
  --deploy-hook-url=env:RENDER_DEPLOY_HOOK_URL
```

Deploy to both Fly.io and Render:

```bash
dagger call --mod ./mkdocs-ci lint-build-publish \
  --ghcr-token=cmd:"gh auth token | tr -d '\n'" \
  --deploy-hook-url=env:RENDER_DEPLOY_HOOK_URL \
  --flyio-app=my-mkdocs-app \
  --flyio-token=env:FLY_API_TOKEN \
  --flyio-region=arn
```

This will:

1. Send a notification to the `athame` topic that the pipeline is starting
2. Run vale, prettier, markdownlint-cli2, and lychee concurrently
3. Send a notification when tests are done (pass or fail)
4. If all tests pass, build the site and publish to GHCR
5. Send a notification when publishing is complete with the image address
6. Trigger Render deploy hook (if provided) and send notification on success/failure
7. Deploy to Fly.io (if app name and token provided) and send notification on success/failure

## Example Output

```
ghcr.io/staticaland/athame/mkdocs-demo:latest@sha256:71a39fda21affdfecd6de5f67b668385578bbc7b9175e3fe0b7558bda83a6275
```

You can then run the container locally:

```bash
docker run -p 8080:80 ghcr.io/staticaland/athame/mkdocs-demo:latest@sha256:...
```

Visit http://localhost:8080 in your browser to view the site.
