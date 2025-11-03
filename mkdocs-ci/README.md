# MkDocs CI

A Dagger module for MkDocs CI/CD: linting, building, and publishing documentation sites.

This module provides functions to lint markdown with vale, prettier, and markdownlint-cli2, build MkDocs Material sites, and publish the site as a container image to ttl.sh registry.

## Features

- **Linting**: Run vale, prettier, and markdownlint-cli2 on your markdown files
- **Concurrent testing**: All linters run in parallel for fast feedback
- **Build**: Build MkDocs Material sites
- **Publish**: Publish sites as container images to ttl.sh
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
```

### Run all tests concurrently

```bash
dagger call --mod ./mkdocs-ci run-all-tests
```

### Build the site

```bash
dagger call --mod ./mkdocs-ci build export --path=./output
```

### Publish to ttl.sh

```bash
dagger call --mod ./mkdocs-ci publish
```

This will:

1. Build the MkDocs Material site
2. Create a container image with nginx and the static files
3. Publish to ttl.sh (available for 1 hour)

### Full CI/CD pipeline

Run all linters, build, and publish if tests pass:

```bash
dagger call --mod ./mkdocs-ci lint-build-publish
```

This will:

1. Send a notification to the `athame` topic that the pipeline is starting
2. Run vale, prettier, and markdownlint-cli2 concurrently
3. Send a notification when tests are done (pass or fail)
4. If all tests pass, build the site and publish to ttl.sh
5. Send a notification when publishing is complete with the image address

## Example Output

```
ttl.sh/mkdocs-demo:1h@sha256:71a39fda21affdfecd6de5f67b668385578bbc7b9175e3fe0b7558bda83a6275
```

You can then run the container locally:

```bash
docker run -p 8080:80 ttl.sh/mkdocs-demo:1h@sha256:...
```

And visit http://localhost:8080 in your browser.
