# MkDocs CI Pipeline

This is a CI/CD pipeline, not a library. Orchestration code implementing continuous delivery principles.

## Pipeline Flow

```
LintBuildPublish → RunAllTests → Build → Publish → deployToAllPlatforms
     └─sequential─┘   └─concurrent─┘         └────concurrent────┘
```

Sequential stages. Concurrent tasks within stages. Uses `errgroup` throughout.

## Key CD Principles

**Fail Fast:** Tests/build/publish stop entire pipeline on first error. Don't waste resources on broken artifacts.

**Build Once, Deploy Many:** `Publish()` creates image with digest (`ghcr.io/user/app@sha256:...`). Same bits deploy everywhere.

**Fast Feedback:** Parallelize independent tasks. All linters run simultaneously. All platforms deploy simultaneously.

**Partial Success (Deploy Only):** Deploy stage runs all platforms even if one fails. Platform-specific failures (credentials, quotas) shouldn't block other platforms.

**No Retries:** Pipeline never retries. Automatic retries hide flaky tests and infrastructure problems. Fix root cause, re-run manually.

## Concurrency Pattern

`errgroup` for parallel execution:

```go
eg, gctx := errgroup.WithContext(ctx)
eg.Go(func() error { return dag.Vale().Check(...) })
eg.Go(func() error { return dag.Prettier().Check(...) })
return eg.Wait()  // First error cancels context, returns immediately
```

## Common Tasks

**Add linter:**
```go
// In RunAllTests()
eg.Go(func() error {
    _, err := dag.NewLinter().Check(...)
    return err
})
```

**Add deployment platform:**
```go
// In deployToAllPlatforms() + add constructor param for credentials
if newPlatformCreds != nil {
    eg.Go(func() error {
        _, err := dag.NewPlatform().Deploy(gctx, ...)
        if err != nil {
            m.notify(gctx, "Failed", opts)
            return fmt.Errorf("deploy failed: %w", err)
        }
        m.notify(gctx, "Succeeded", opts)
        return nil
    })
}
```

**Test stages individually:**
```bash
dagger call --mod ./mkdocs-ci run-all-tests
dagger call --mod ./mkdocs-ci build export --path ./output
dagger call --mod ./mkdocs-ci publish --ghcr-token=env:GITHUB_TOKEN
```

## Multi-Platform Builds

Publishes `linux/amd64` (Render requirement) and `linux/arm64` (Apple Silicon). Append to `platforms` slice in `Publish()` to add more.

## Notifications

`notify()` never blocks pipeline. Logs errors but doesn't return them. Makes pipeline observable without affecting reliability.
