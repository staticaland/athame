# CI Pipeline Implementation Plan: Verify/Publish/Deploy

## Overview

Restructure `mkdocs-ci` and `miele-ci` into three infrastructure-based phases:

- **Verify** - Local validation (no credentials needed): lint → build → test → scan
- **Publish** - Push to registry (requires GHCR token)
- **Deploy** - Deploy to cloud (requires cloud credentials)

## Design Principles

1. **Verify returns the validated artifact** - enables reuse in Publish without rebuilding
2. **Security scans before publish** - never push vulnerable images to registry
3. **Local-first DX** - developers can run Verify without any secrets
4. **Composable** - each phase can be called independently or chained

## Phase Boundaries

### What Requires Infrastructure?

| Phase | Needs | Can Run Locally? |
|-------|-------|------------------|
| Verify | None | ✅ Yes |
| Publish | Registry credentials | ❌ No |
| Deploy | Cloud credentials, running services | ❌ No |

## Implementation: miele-ci

### Current State

```go
// Current monolithic function
func (m *MieleCi) BuildPublishDeploy(
    ctx context.Context,
    ghcrToken *dagger.Secret,
    flyioApp string,
    flyioToken *dagger.Secret,
    flyioRegion string,
) (string, error)
```

### Target State

```go
// Phase 1: Verify (local-friendly)
// Returns the verified, scanned container ready for publish
func (m *MieleCi) Verify(ctx context.Context) (*dagger.Container, error) {
    // 1. Lint (future: add ESLint/Prettier)
    // 2. Test
    output, err := m.Test(ctx)
    if err != nil {
        return nil, fmt.Errorf("tests failed: %w", err)
    }
    fmt.Printf("Tests passed:\n%s\n", output)

    // 3. Build artifact
    builtSite := m.Build()

    // 4. Create multi-platform container variants
    platforms := []dagger.Platform{
        "linux/amd64",
        "linux/arm64",
    }

    platformVariants := make([]*dagger.Container, 0, len(platforms))
    for _, platform := range platforms {
        ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
            // renovate: datasource=docker depName=nginx
            From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
            WithDirectory("/usr/share/nginx/html", builtSite).
            WithExposedPort(80).
            WithLabel("org.opencontainers.image.title", m.ImageName).
            WithLabel("org.opencontainers.image.version", m.Tag).
            WithLabel("org.opencontainers.image.created", time.Now().String()).
            WithLabel("org.opencontainers.image.source", "https://github.com/staticaland/athame")

        platformVariants = append(platformVariants, ctr)
    }

    // 5. Security scan the first platform variant (representative)
    m.notify(ctx, "Scanning container for vulnerabilities...", dagger.NtfySendOpts{
        Title:    "Trivy Security Scan Started",
        Priority: "default",
        Tags:     "shield",
    })

    scanResult, err := dag.Trivy().ScanContainer(ctx, platformVariants[0], "scan-target")
    if err != nil {
        m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
            Title:    "Trivy Security Scan Failed",
            Priority: "high",
            Tags:     "warning",
        })
        return nil, fmt.Errorf("trivy scan failed: %w", err)
    }
    fmt.Printf("Trivy scan results:\n%s\n", scanResult)

    m.notify(ctx, "Security scan completed successfully.", dagger.NtfySendOpts{
        Title:    "Trivy Security Scan Completed",
        Priority: "default",
        Tags:     "white_check_mark",
    })

    // Return the verified container (first platform as representative)
    // Note: For multi-platform, we'd need to store all variants
    return platformVariants[0], nil
}

// Phase 2: Publish (requires registry credentials)
// Takes verified container and publishes to GHCR
func (m *MieleCi) Publish(
    ctx context.Context,
    // GitHub token for GHCR authentication (get with: gh auth token)
    ghcrToken *dagger.Secret,
) (string, error) {
    // Get verified container by running Verify again
    // Dagger's caching ensures this doesn't rebuild
    _, err := m.Verify(ctx)
    if err != nil {
        return "", fmt.Errorf("verification failed: %w", err)
    }

    // Rebuild platform variants (cached from Verify)
    builtSite := m.Build()
    platforms := []dagger.Platform{
        "linux/amd64",
        "linux/arm64",
    }

    platformVariants := make([]*dagger.Container, 0, len(platforms))
    for _, platform := range platforms {
        ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
            From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
            WithDirectory("/usr/share/nginx/html", builtSite).
            WithExposedPort(80).
            WithLabel("org.opencontainers.image.title", m.ImageName).
            WithLabel("org.opencontainers.image.version", m.Tag).
            WithLabel("org.opencontainers.image.created", time.Now().String()).
            WithLabel("org.opencontainers.image.source", "https://github.com/staticaland/athame")

        platformVariants = append(platformVariants, ctr)
    }

    // Publish to GHCR
    imageAddr := fmt.Sprintf("ghcr.io/%s/athame/%s:%s", m.GhcrUsername, m.ImageName, m.Tag)
    addr, err := dag.Container().
        WithRegistryAuth("ghcr.io", m.GhcrUsername, ghcrToken).
        Publish(ctx, imageAddr, dagger.ContainerPublishOpts{
            PlatformVariants: platformVariants,
        })
    if err != nil {
        return "", fmt.Errorf("failed to publish to GHCR: %w", err)
    }

    m.notify(ctx,
        fmt.Sprintf("Published to GHCR.\n\n**Image:**\n```\n%s\n```", addr),
        dagger.NtfySendOpts{
            Title:    "Image Publishing Completed",
            Priority: "default",
            Tags:     "white_check_mark",
            Markdown: true,
        })

    return addr, nil
}

// Phase 3: Deploy (requires cloud credentials)
// Takes published image address and deploys to Fly.io
func (m *MieleCi) Deploy(
    ctx context.Context,
    // Image address from Publish (e.g., ghcr.io/user/image:tag@sha256:...)
    imageAddr string,
    // Fly.io app name
    flyioApp string,
    // Fly.io API token
    flyioToken *dagger.Secret,
    // +optional
    // +default="arn"
    flyioRegion string,
) error {
    _, err := dag.Flyio().Deploy(ctx, flyioApp, imageAddr, flyioToken, dagger.FlyioDeployOpts{
        PrimaryRegion: flyioRegion,
        InternalPort:  80,
    })
    if err != nil {
        m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
            Title:    "Fly.io Deploy Failed",
            Priority: "high",
            Tags:     "warning",
        })
        return fmt.Errorf("fly.io deploy failed: %w", err)
    }

    flyioUrl := fmt.Sprintf("https://%s.fly.dev", flyioApp)
    m.notify(ctx,
        fmt.Sprintf("Deployed to Fly.io.\n\n**App:** %s", flyioApp),
        dagger.NtfySendOpts{
            Title:    "Fly.io Deploy Completed",
            Priority: "default",
            Tags:     "white_check_mark",
            Actions:  fmt.Sprintf("view, View Site, %s", flyioUrl),
            Markdown: true,
        })

    return nil
}

// Convenience: Full pipeline
// Chains all three phases for CI/CD workflows
func (m *MieleCi) Pipeline(
    ctx context.Context,
    ghcrToken *dagger.Secret,
    flyioApp string,
    flyioToken *dagger.Secret,
    // +optional
    // +default="arn"
    flyioRegion string,
) (string, error) {
    m.notify(ctx, "Starting pipeline...", dagger.NtfySendOpts{
        Title:    "Miele CI/CD Started",
        Priority: "default",
        Tags:     "hourglass_flowing_sand",
    })

    // Phase 1: Verify
    _, err := m.Verify(ctx)
    if err != nil {
        return "", err
    }

    // Phase 2: Publish
    addr, err := m.Publish(ctx, ghcrToken)
    if err != nil {
        return "", err
    }

    // Phase 3: Deploy
    err = m.Deploy(ctx, addr, flyioApp, flyioToken, flyioRegion)
    if err != nil {
        return addr, err // Return addr even on deploy failure
    }

    return addr, nil
}
```

### Developer Workflows

#### Local Development (No Secrets)
```bash
# Verify everything before pushing
dagger call --mod ./miele-ci verify

# Output:
# ✅ Tests passed
# ✅ Build succeeded
# ✅ Container built
# ✅ Trivy scan: no vulnerabilities
# 🎉 Safe to push!
```

#### CI/CD (GitHub Actions)
```yaml
- name: Verify
  run: dagger call --mod ./miele-ci verify

- name: Publish
  run: |
    ADDR=$(dagger call --mod ./miele-ci publish \
      --ghcr-token=env:GITHUB_TOKEN)
    echo "IMAGE_ADDR=$ADDR" >> $GITHUB_OUTPUT

- name: Deploy
  run: |
    dagger call --mod ./miele-ci deploy \
      --image-addr="${{ steps.publish.outputs.IMAGE_ADDR }}" \
      --flyio-app=my-app \
      --flyio-token=env:FLY_TOKEN
```

#### Manual Publish/Deploy (Maintainer)
```bash
# Verify first
dagger call --mod ./miele-ci verify

# Publish to registry
ADDR=$(dagger call --mod ./miele-ci publish --ghcr-token=env:GITHUB_TOKEN)

# Deploy to Fly.io
dagger call --mod ./miele-ci deploy \
  --image-addr="$ADDR" \
  --flyio-app=miele \
  --flyio-token=env:FLY_TOKEN
```

## Implementation: mkdocs-ci

### Current State

```go
func (m *MkdocsCi) TestBuildPublishDeploy(
    ctx context.Context,
    ghcrToken *dagger.Secret,
    deployHookURL *dagger.Secret,
    flyioApp string,
    flyioToken *dagger.Secret,
    // ... many more parameters
) (string, error)
```

### Target State

```go
// Phase 1: Verify
func (m *MkdocsCi) Verify(ctx context.Context) (*dagger.Container, error) {
    siteDir := m.Source.Directory(m.SitePath)

    // 1. Lint in parallel
    eg, gctx := errgroup.WithContext(ctx)

    eg.Go(func() error {
        _, err := dag.Vale().Check(dagger.ValeCheckOpts{
            Source: siteDir,
            Path:   "docs",
        }).Stdout(gctx)
        return err
    })

    eg.Go(func() error {
        _, err := dag.Prettier().Check(dagger.PrettierCheckOpts{
            Source:  siteDir,
            Pattern: "docs/**/*.md",
        }).Stdout(gctx)
        return err
    })

    eg.Go(func() error {
        _, err := dag.MarkdownlintCli2().Check(dagger.MarkdownlintCli2CheckOpts{
            Source:  siteDir,
            Pattern: "docs/**/*.md",
        }).Stdout(gctx)
        return err
    })

    if err := eg.Wait(); err != nil {
        m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
            Title:    "Linting Failed",
            Priority: "high",
            Tags:     "warning",
        })
        return nil, fmt.Errorf("linting failed: %w", err)
    }

    m.notify(ctx, "Linting passed. Building site...", dagger.NtfySendOpts{
        Title:    "Linting Completed",
        Priority: "default",
        Tags:     "white_check_mark",
    })

    // 2. Build
    builtSite := m.Build()

    // 3. Link checking (runs on built site)
    _, err := dag.Lychee().Check(dagger.LycheeCheckOpts{
        Source: builtSite,
        Path:   ".",
    }).Stdout(ctx)
    if err != nil {
        m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
            Title:    "Link Checking Failed",
            Priority: "high",
            Tags:     "warning",
        })
        return nil, fmt.Errorf("link checking failed: %w", err)
    }

    m.notify(ctx, "Link checking passed. Building container...", dagger.NtfySendOpts{
        Title:    "Link Checking Completed",
        Priority: "default",
        Tags:     "white_check_mark",
    })

    // 4. Create multi-platform container
    platforms := []dagger.Platform{
        "linux/amd64",
        "linux/arm64",
    }

    platformVariants := make([]*dagger.Container, 0, len(platforms))
    for _, platform := range platforms {
        ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
            // renovate: datasource=docker depName=nginx
            From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
            WithDirectory("/usr/share/nginx/html", builtSite).
            WithExposedPort(80).
            WithLabel("org.opencontainers.image.title", m.ImageName).
            WithLabel("org.opencontainers.image.version", m.Tag).
            WithLabel("org.opencontainers.image.created", time.Now().String()).
            WithLabel("org.opencontainers.image.source", "https://github.com/staticaland/athame")

        platformVariants = append(platformVariants, ctr)
    }

    // 5. Security scan (optional for mkdocs, but good practice)
    scanResult, err := dag.Trivy().ScanContainer(ctx, platformVariants[0], "scan-target")
    if err != nil {
        m.notify(ctx, "Check logs for details.", dagger.NtfySendOpts{
            Title:    "Security Scan Failed",
            Priority: "high",
            Tags:     "warning",
        })
        return nil, fmt.Errorf("security scan failed: %w", err)
    }
    fmt.Printf("Trivy scan results:\n%s\n", scanResult)

    m.notify(ctx, "All verification passed. Ready to publish.", dagger.NtfySendOpts{
        Title:    "Verification Completed",
        Priority: "default",
        Tags:     "white_check_mark",
    })

    return platformVariants[0], nil
}

// Phase 2: Publish
func (m *MkdocsCi) Publish(
    ctx context.Context,
    ghcrToken *dagger.Secret,
) (string, error) {
    // Verify first (cached by Dagger)
    _, err := m.Verify(ctx)
    if err != nil {
        return "", err
    }

    // Rebuild container (cached)
    builtSite := m.Build()
    platforms := []dagger.Platform{
        "linux/amd64",
        "linux/arm64",
    }

    platformVariants := make([]*dagger.Container, 0, len(platforms))
    for _, platform := range platforms {
        ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
            From("nginx:1.27.5-alpine3.21@sha256:65645c7bb6a0661892a8b03b89d0743208a18dd2f3f17a54ef4b76fb8e2f2a10").
            WithDirectory("/usr/share/nginx/html", builtSite).
            WithExposedPort(80).
            WithLabel("org.opencontainers.image.title", m.ImageName).
            WithLabel("org.opencontainers.image.version", m.Tag).
            WithLabel("org.opencontainers.image.created", time.Now().String()).
            WithLabel("org.opencontainers.image.source", "https://github.com/staticaland/athame")

        platformVariants = append(platformVariants, ctr)
    }

    // Publish to GHCR
    imageAddr := fmt.Sprintf("ghcr.io/%s/athame/%s:%s", m.GhcrUsername, m.ImageName, m.Tag)
    addr, err := dag.Container().
        WithRegistryAuth("ghcr.io", m.GhcrUsername, ghcrToken).
        Publish(ctx, imageAddr, dagger.ContainerPublishOpts{
            PlatformVariants: platformVariants,
        })
    if err != nil {
        return "", fmt.Errorf("failed to publish to GHCR: %w", err)
    }

    m.notify(ctx,
        fmt.Sprintf("Published to GHCR.\n\n**Image:**\n```\n%s\n```", addr),
        dagger.NtfySendOpts{
            Title:    "Image Publishing Completed",
            Priority: "default",
            Tags:     "white_check_mark",
            Markdown: true,
        })

    return addr, nil
}

// Phase 3: Deploy
func (m *MkdocsCi) Deploy(
    ctx context.Context,
    imageAddr string,
    // +optional
    deployHookURL *dagger.Secret,
    // +optional
    flyioApp string,
    // +optional
    flyioToken *dagger.Secret,
    // +optional
    // +default="arn"
    flyioRegion string,
    // +optional
    gcloudService string,
    // +optional
    gcloudProject string,
    // +optional
    gcloudRegion string,
    // +optional
    gcloudServiceAccountKey *dagger.Secret,
    // +optional
    gcloudAllowUnauthenticated bool,
    // +optional
    // +default="ghcr"
    artifactRegistryRepo string,
    // +optional
    // +default="europe-north2"
    artifactRegistryRegion string,
) error {
    return m.deployToAllPlatforms(
        ctx,
        imageAddr,
        deployHookURL,
        flyioApp,
        flyioToken,
        flyioRegion,
        gcloudService,
        gcloudProject,
        gcloudRegion,
        gcloudServiceAccountKey,
        gcloudAllowUnauthenticated,
        artifactRegistryRepo,
        artifactRegistryRegion,
    )
}

// Convenience: Full pipeline
func (m *MkdocsCi) Pipeline(
    ctx context.Context,
    ghcrToken *dagger.Secret,
    // ... all deploy parameters
) (string, error) {
    m.notify(ctx, "Starting pipeline...", dagger.NtfySendOpts{
        Title:    "MkDocs CI/CD Started",
        Priority: "default",
        Tags:     "hourglass_flowing_sand",
    })

    _, err := m.Verify(ctx)
    if err != nil {
        return "", err
    }

    addr, err := m.Publish(ctx, ghcrToken)
    if err != nil {
        return "", err
    }

    err = m.Deploy(ctx, addr, /* ... deploy params */)
    return addr, err
}
```

### Developer Workflows

#### Local Development
```bash
# Verify everything before pushing
dagger call --mod ./mkdocs-ci verify

# Output:
# ✅ Vale: no issues
# ✅ Prettier: formatted correctly
# ✅ Markdownlint: no issues
# ✅ Build succeeded
# ✅ Lychee: all links valid
# ✅ Trivy: no vulnerabilities
# 🎉 Safe to push!
```

#### CI/CD with Multiple Deployments
```yaml
- name: Verify
  run: dagger call --mod ./mkdocs-ci verify

- name: Publish
  run: |
    ADDR=$(dagger call --mod ./mkdocs-ci publish \
      --ghcr-token=env:GITHUB_TOKEN)
    echo "IMAGE_ADDR=$ADDR" >> $GITHUB_OUTPUT

- name: Deploy to Render
  run: |
    dagger call --mod ./mkdocs-ci deploy \
      --image-addr="${{ steps.publish.outputs.IMAGE_ADDR }}" \
      --deploy-hook-url=env:RENDER_DEPLOY_HOOK

- name: Deploy to Fly.io
  run: |
    dagger call --mod ./mkdocs-ci deploy \
      --image-addr="${{ steps.publish.outputs.IMAGE_ADDR }}" \
      --flyio-app=mkdocs-demo \
      --flyio-token=env:FLY_TOKEN

- name: Deploy to Cloud Run
  run: |
    dagger call --mod ./mkdocs-ci deploy \
      --image-addr="${{ steps.publish.outputs.IMAGE_ADDR }}" \
      --gcloud-service=mkdocs-demo \
      --gcloud-project=my-project \
      --gcloud-service-account-key=env:GCLOUD_KEY
```

## Migration Strategy

### Phase 1: Add New Functions (Non-Breaking)
1. Add `Verify()`, `Publish()`, `Deploy()` to both modules
2. Keep existing `TestBuildPublishDeploy()` / `BuildPublishDeploy()`
3. Update existing functions to call new phases internally
4. Test both paths work identically

### Phase 2: Update Documentation
1. Update README to show new workflow
2. Add examples for local verification
3. Document CI/CD patterns

### Phase 3: Deprecate Old Functions
1. Mark old functions as deprecated in comments
2. Update all internal usage to new phases
3. Monitor for external usage

### Phase 4: Remove Old Functions (Breaking)
1. Remove deprecated functions in next major version
2. Update CHANGELOG with migration guide

## Benefits Summary

### For Developers
- ✅ Run full validation locally without secrets
- ✅ Know if push will succeed before pushing
- ✅ Fast feedback loop
- ✅ No need to set up registry credentials locally

### For CI/CD
- ✅ Clear separation of concerns
- ✅ Easy to parallelize independent deployments
- ✅ Can publish once, deploy many times
- ✅ Better error messages (know which phase failed)

### For Security
- ✅ Never publish vulnerable images
- ✅ Scan happens before any credentials are needed
- ✅ Audit trail: verified → published → deployed

## Open Questions

1. **Container storage between phases**: Currently rebuilding in Publish (relying on cache). Should we store the verified container artifact differently?

2. **Parallel deployments**: Should Deploy support multiple targets in one call, or prefer multiple Deploy calls?

3. **Naming**: Keep `BuildPublishDeploy` as alias to `Pipeline` for backward compatibility?

4. **Test function**: Should `Test()` be separate or part of `Verify()`? Current plan includes it in Verify.

5. **Notifications**: Should notifications be in each phase or only in Pipeline?
