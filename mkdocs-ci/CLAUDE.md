# MkDocs CI/CD Pipeline

## What This Is

This is a **CI/CD pipeline**, not a library. It orchestrates the complete build-test-deploy workflow for MkDocs Material sites.

## Pipeline Mindset

Think like a pipeline engineer:

- **Sequential stages** - tests → build → publish → deploy
- **Fast failure** - stop at first error, don't waste time on downstream stages
- **Parallel execution** - run independent tasks concurrently (all linters, all deployments)
- **Portability** - same pipeline runs locally and in CI (GitHub Actions, GitLab CI, etc.)

## Pipeline Stages

1. **Test** - runs vale, prettier, markdownlint, lychee concurrently
2. **Build** - generates static site with MkDocs Material
3. **Publish** - packages as multi-platform container image to GHCR
4. **Deploy** - ships to Render, Fly.io, and Google Cloud Run concurrently

## Run Anywhere

**Local:**
```bash
dagger call --mod ./mkdocs-ci lint-build-publish --ghcr-token=env:GITHUB_TOKEN
```

**GitHub Actions:**
```yaml
- name: Deploy
  run: dagger call --mod ./mkdocs-ci lint-build-publish --ghcr-token=env:GITHUB_TOKEN
```

Same command. Same result.

## Configuration

Pipeline configured via constructor parameters:

- `source` - repository root
- `sitePath` - path to MkDocs project within repo
- `imageName` - container image name
- `ghcrUsername` - GitHub Container Registry username
- Deployment credentials (render, flyio, gcloud) - passed to `LintBuildPublish`

## Exit Codes

- **0** - all stages succeeded
- **Non-zero** - stage failed, check logs

No retries. Fix the issue and re-run.
