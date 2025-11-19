# MkDocs CI Pipeline

This is a CI/CD pipeline, not a library. You're working on orchestration code implementing continuous delivery principles.

## Continuous Delivery Principles

This pipeline follows patterns from *Continuous Delivery* (Humble & Farley):

### 1. Fail Fast
**Principle:** Stop at the first sign of failure. Don't waste resources on downstream stages.

**Implementation:** Sequential stages with early exit:
```
tests fail → pipeline stops (no build, no publish, no deploy)
build fails → pipeline stops (no publish, no deploy)
publish fails → pipeline stops (no deploy)
```

Each stage blocks the next. Tests run concurrently, but pipeline waits for all before proceeding.

### 2. Build Once, Deploy Many
**Principle:** Build your artifact once, deploy the same artifact to all environments.

**Implementation:** `Publish()` creates multi-platform container image with digest. Same image deployed to all platforms:
```go
addr, err := m.Publish(ctx, ghcrToken)  // Build once: ghcr.io/user/app:tag@sha256:...
// Deploy same digest to Render, Fly.io, Cloud Run
deployToAllPlatforms(ctx, addr, ...)
```

The `@sha256:...` digest ensures bit-for-bit identical artifact everywhere.

### 3. Optimize for Fast Feedback
**Principle:** Keep the pipeline fast. Parallelize independent tasks.

**Implementation:** `errgroup` for concurrent execution within stages:
- **Tests:** All linters run simultaneously (vale, prettier, markdownlint, lychee)
- **Deploys:** All platforms deploy simultaneously (Render, Fly.io, Cloud Run)
- **Builds:** Multi-platform container builds run in parallel

### 4. Repeatable and Reliable
**Principle:** Same inputs → same outputs. Local == CI.

**Implementation:** Dagger containers provide hermetic builds. Same command works anywhere:
```bash
# On your laptop
dagger call --mod ./mkdocs-ci lint-build-publish --ghcr-token=env:GITHUB_TOKEN

# In GitHub Actions
dagger call --mod ./mkdocs-ci lint-build-publish --ghcr-token=env:GITHUB_TOKEN
```

No "works on my machine" problems.

### 5. Make the Process Observable
**Principle:** Pipeline should broadcast its state. Silence is not golden.

**Implementation:** `notify()` sends notifications at every stage transition:
- Tests started
- Tests passed/failed
- Build complete
- Each platform deployed
- Overall success/failure

Notifications never block the pipeline (non-blocking, log errors only).

## Code Architecture

**Pipeline flow:**
```
LintBuildPublish → RunAllTests → Build → Publish → deployToAllPlatforms
     └─sequential─┘   └─concurrent─┘         └────concurrent────┘
```

**Core pattern:** `errgroup` for concurrent execution with fail-fast behavior.

## Implementation Details

### Concurrency Pattern

Uses `errgroup` from `golang.org/x/sync/errgroup`:

```go
eg, gctx := errgroup.WithContext(ctx)
eg.Go(func() error { /* task 1 */ })
eg.Go(func() error { /* task 2 */ })
return eg.Wait()  // Wait for all, return first error
```

**Why errgroup:**
- Automatic cancellation: First error cancels context for other goroutines
- Simplified error handling: Returns first error encountered
- Race-free: Safe concurrent execution without manual synchronization

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

**Why all platforms complete:** Deploy stage uses same `errgroup` pattern, but each goroutine handles its own errors and notifications. Even if Render fails, Fly.io and Cloud Run still attempt deployment. This provides partial success - at least some environments get updated.

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

## Error Handling Philosophy

### Stop the Line (Tests/Build/Publish)
**CD Principle:** Treat every failure as an opportunity to improve. Don't bypass quality gates.

Early stages (test, build, publish) implement strict fail-fast:
- Single linter failure stops entire pipeline
- Build failure prevents publishing
- Publish failure prevents deployment

**Why:** These stages validate quality. Bypassing them means deploying broken artifacts.

### Partial Success (Deploy)
**CD Principle:** Maximize availability. Some environments updated is better than none.

Deploy stage attempts all platforms even if one fails:
- Render fails → Fly.io and Cloud Run still deploy
- Returns first error but completes all deployments
- Notifications indicate which platforms succeeded/failed

**Why:** Deployment failures are often platform-specific (credentials, quotas). Other platforms may succeed.

### No Automatic Retries
**CD Principle:** Retries mask problems. Fix the root cause instead.

Pipeline never retries automatically:
- Tests fail → fix the code, re-run pipeline
- Deploy fails → check credentials, re-run pipeline
- Network glitch → re-run pipeline manually

**Why:** Automatic retries hide flaky tests, intermittent failures, and infrastructure problems. Manual re-run forces acknowledgment of the failure.

## Common Modifications

### Add a Quality Gate (Linter)
Add `eg.Go()` block in `RunAllTests()`:
```go
eg.Go(func() error {
    _, err := dag.NewLinterModule().Check(...)
    return err
})
```

**CD consideration:** Each linter is a quality gate. Add linters that catch real problems, not style preferences. Keep the test stage fast (< 5 minutes ideal).

### Add a Deployment Platform
Add conditional `eg.Go()` block in `deployToAllPlatforms()` + add constructor params for credentials:
```go
if newPlatformCreds != nil {
    eg.Go(func() error {
        _, err := dag.NewPlatform().Deploy(gctx, ...)
        if err != nil {
            m.notify(gctx, "Check logs", opts) // always notify
            return fmt.Errorf("platform deploy failed: %w", err)
        }
        m.notify(gctx, "Deploy succeeded", opts) // success notification
        return nil
    })
}
```

**CD consideration:** New platform deploys concurrently with existing platforms. No impact on deploy time (unless you hit parallel execution limits).

### Change Base Image
Update renovate comment + digest in `Publish()` (currently nginx):
```go
// renovate: datasource=docker depName=nginx
From("nginx:1.27.5-alpine3.21@sha256:...")
```

**CD consideration:** Base images are part of your artifact. Pin with digest (`@sha256:...`) for reproducible builds. Let Renovate handle updates.

### Add/Modify Notifications
Update `notify()` calls in `LintBuildPublish()` and `deployToAllPlatforms()`:
```go
m.notify(ctx, "Your message", dagger.NtfySendOpts{
    Title:    "Stage Complete",
    Priority: "default",  // or "high" for failures
    Tags:     "white_check_mark",  // or "warning" for failures
})
```

**CD consideration:** Notifications make the pipeline observable. Add notifications for state transitions and failures. Never let notifications block the pipeline (they log errors but don't return them).
