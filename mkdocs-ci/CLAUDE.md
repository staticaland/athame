# MkDocs CI Pipeline

This is a CI/CD pipeline, not a library. You're working on orchestration code that runs stages sequentially but executes tasks within each stage concurrently.

## Code Architecture

**Pipeline flow:**
```
LintBuildPublish → RunAllTests → Build → Publish → deployToAllPlatforms
     └─sequential─┘   └─concurrent─┘         └────concurrent────┘
```

**Core pattern:** `errgroup` for concurrent execution with fail-fast behavior.

## Concurrency Model

### Test Stage (`RunAllTests`)

All linters run concurrently. First failure stops the pipeline:

```go
eg.Go(func() error {
    _, err := dag.Vale().Check(...)
    return err
})
eg.Go(func() error {
    _, err := dag.Prettier().Check(...)
    return err
})
// ... more linters
return eg.Wait() // blocks until all complete or first error
```

**Adding a new linter:** Add another `eg.Go()` block. No coordination needed.

### Deploy Stage (`deployToAllPlatforms`)

All platforms deploy concurrently. All run to completion even if one fails:

```go
if deployHookURL != nil {
    eg.Go(func() error {
        _, err := dag.RenderDeployHook(deployHookURL).Deploy(gctx)
        if err != nil {
            m.notify(gctx, ...) // notify but continue
            return fmt.Errorf("render deploy failed: %w", err)
        }
        return nil
    })
}
// ... more platforms
return eg.Wait() // returns first error but all goroutines complete
```

**Adding a new platform:** Add another conditional `eg.Go()` block. Pattern: check credentials exist, deploy, notify on success/failure.

## Notifications

`notify()` sends ntfy notifications. Always non-blocking - logs errors but never fails the pipeline:

```go
func (m *MkdocsCi) notify(ctx context.Context, message string, opts dagger.NtfySendOpts) {
    _, err := dag.Ntfy().Send(ctx, "athame", message, opts)
    if err != nil {
        fmt.Printf("Failed to send notification '%s': %v\n", opts.Title, err)
    }
}
```

**Call sites:** stage transitions (tests started, tests passed, build complete, etc.) and deployment outcomes.

## Testing Changes

**Test individual stages:**
```bash
# Just run tests
dagger call --mod ./mkdocs-ci run-all-tests

# Just build
dagger call --mod ./mkdocs-ci build export --path ./output

# Build and publish (requires GHCR token)
dagger call --mod ./mkdocs-ci publish --ghcr-token=env:GITHUB_TOKEN
```

**Full pipeline:**
```bash
dagger call --mod ./mkdocs-ci lint-build-publish \
    --ghcr-token=env:GITHUB_TOKEN \
    --deploy-hook-url=env:RENDER_DEPLOY_HOOK
```

## Multi-Platform Images

Publishes `linux/amd64` (required for Render) and `linux/arm64` (for Apple Silicon):

```go
platforms := []dagger.Platform{
    "linux/amd64",
    "linux/arm64",
}
for _, platform := range platforms {
    ctr := dag.Container(dagger.ContainerOpts{Platform: platform}).
        From("nginx:...@sha256:...").
        WithDirectory("/usr/share/nginx/html", builtSite)
    platformVariants = append(platformVariants, ctr)
}
```

**Adding platforms:** Append to `platforms` slice. Increases publish time linearly.

## Error Handling

**Fast failure:** Tests stage stops at first error. Subsequent stages don't run.

**Partial success:** Deploy stage runs all platforms. Returns first error but attempts all deploys.

**No retries:** Pipeline doesn't retry. Re-run the entire pipeline on failure.

## Common Modifications

**Add a linter:** Add `eg.Go()` block in `RunAllTests()`

**Add a deployment platform:** Add conditional `eg.Go()` block in `deployToAllPlatforms()` + add constructor params for credentials

**Change base image:** Update renovate comment + digest in `Publish()` (currently nginx)

**Modify notifications:** Update `notify()` calls in `LintBuildPublish()` and `deployToAllPlatforms()`
