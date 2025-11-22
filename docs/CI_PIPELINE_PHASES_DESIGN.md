# CI/CD Pipeline Phases Design

## Current State Analysis

### mkdocs-ci Functions
- **Test**: Vale, Prettier, Markdownlint, Lychee (concurrent)
- **Build**: MkDocs Material site
- **BuildPublish**: Multi-platform container to GHCR
- **TestBuildPublishDeploy**: Full pipeline to Render/Fly.io/Google Cloud Run

### miele-ci Functions
- **Test**: npm test
- **Build**: Vite application
- **BuildPublish**: Multi-platform container to GHCR with Trivy scan
- **BuildPublishDeploy**: Full pipeline to Fly.io

---

## Pipeline Phase Philosophies

Based on Dave Farley's Continuous Delivery principles, pipeline phases should:

1. **Provide fast feedback** - fail fast on simple problems
2. **Progress from cheap to expensive** - run quick checks before slow ones
3. **Progress from common to rare failures** - catch typical errors early
4. **Be independently deployable** - each stage is a quality gate
5. **Maintain reproducibility** - same input = same output

---

## Design Option 1: Classic Three-Stage Pipeline

**Philosophy**: Traditional "Integration → Delivery → Deployment" with clear separation of concerns.

### Phases

#### Integration (Commit Stage)
Fast feedback on code quality - runs on every commit.

```
mkdocs-ci:
  - Integrate()
    - Lint (vale, prettier, markdownlint)
    - Build verification
    - Fast smoke tests

miele-ci:
  - Integrate()
    - npm test (unit tests)
    - Build verification
    - ESLint/Prettier
```

#### Delivery (Acceptance Stage)
Creates deployable artifacts and validates them.

```
mkdocs-ci:
  - Deliver()
    - Build MkDocs site
    - Link checking (lychee)
    - Publish to GHCR
    - Optional: visual regression tests

miele-ci:
  - Deliver()
    - Build Vite app
    - Security scan (Trivy)
    - Publish to GHCR
    - Optional: E2E tests
```

#### Deployment
Progressively deploys to environments.

```
mkdocs-ci:
  - DeployStaging() → test environment
  - DeployProduction() → prod environment
    - Render, Fly.io, or Cloud Run

miele-ci:
  - DeployStaging() → test environment
  - DeployProduction() → Fly.io
```

### Strengths
- ✅ Clear mental model - classic CI/CD terminology
- ✅ Well-understood by most teams
- ✅ Natural progression: quality → artifact → release
- ✅ Easy to visualize in CI/CD dashboards

### Weaknesses
- ⚠️ "Integration" might be confusing (integrate with what?)
- ⚠️ "Delivery" vs "Deployment" distinction is subtle
- ⚠️ Doesn't emphasize security/quality gates explicitly
- ⚠️ May feel too abstract for newcomers

### When to Use
- Teams familiar with CD terminology
- Organizations with formal release processes
- When you need clear separation between "artifact ready" and "deployed"

---

## Design Option 2: Validation Pipeline (Quality-Focused)

**Philosophy**: Emphasizes increasing levels of validation and confidence.

### Phases

#### Validate-Code
Source code quality checks.

```
mkdocs-ci:
  - ValidateCode()
    - Vale (prose linting)
    - Prettier (formatting)
    - Markdownlint (structure)

miele-ci:
  - ValidateCode()
    - ESLint
    - Prettier
    - TypeScript type checking
```

#### Validate-Build
Verify the build process and artifact.

```
mkdocs-ci:
  - ValidateBuild()
    - Build MkDocs site
    - Lychee link checking
    - Size checks

miele-ci:
  - ValidateBuild()
    - Vite build
    - npm test
    - Bundle analysis
```

#### Validate-Security
Security and compliance checks.

```
mkdocs-ci:
  - ValidateSecurity()
    - Container image scanning
    - Dependency audit
    - SBOM generation

miele-ci:
  - ValidateSecurity()
    - Trivy container scan
    - npm audit
    - SBOM generation
```

#### Release
Publish and deploy verified artifacts.

```
mkdocs-ci:
  - Release()
    - Publish to GHCR
    - Deploy to environment(s)

miele-ci:
  - Release()
    - Publish to GHCR
    - Deploy to Fly.io
```

### Strengths
- ✅ Explicit security gate
- ✅ "Validate" is clear and actionable
- ✅ Each phase has obvious pass/fail criteria
- ✅ Easy to understand what each stage does
- ✅ Scalable - can add more validation phases

### Weaknesses
- ⚠️ More phases = longer pipeline
- ⚠️ May be overkill for simple projects
- ⚠️ "Validate" feels repetitive in names

### When to Use
- Security-conscious organizations
- Projects with compliance requirements
- Teams that want explicit quality gates
- When you need granular failure reporting

---

## Design Option 3: Progressive Confidence Pipeline

**Philosophy**: Build confidence progressively with increasing investment.

### Phases

#### Quick-Check (< 2 minutes)
Immediate feedback on common mistakes.

```
mkdocs-ci:
  - QuickCheck()
    - Prettier check only
    - Basic markdown lint

miele-ci:
  - QuickCheck()
    - Prettier check
    - TypeScript compile check
```

#### Full-Check (< 10 minutes)
Comprehensive pre-build validation.

```
mkdocs-ci:
  - FullCheck()
    - All linters (vale, prettier, markdownlint)
    - Build
    - Link checking

miele-ci:
  - FullCheck()
    - All linters
    - npm test
    - Build
```

#### Package (< 15 minutes)
Create and verify deployable artifact.

```
mkdocs-ci:
  - Package()
    - Multi-platform container build
    - Security scan
    - Publish to GHCR

miele-ci:
  - Package()
    - Multi-platform container build
    - Trivy scan
    - Publish to GHCR
```

#### Ship
Deploy to environments.

```
mkdocs-ci:
  - Ship(environment)
    - Deploy based on branch/tag
    - Smoke tests post-deployment

miele-ci:
  - Ship(environment)
    - Deploy to Fly.io
    - Health check
```

### Strengths
- ✅ Time-based budgets are intuitive
- ✅ Clear optimization target (keep Quick-Check under 2 min)
- ✅ "Ship" is more modern than "Deploy"
- ✅ Simple, memorable phase names

### Weaknesses
- ⚠️ Time budgets may vary by project size
- ⚠️ Less prescriptive about what goes where
- ⚠️ "Check" appears in two phase names

### When to Use
- Developer experience is priority
- Fast feedback is critical
- Teams value simplicity over formality
- When optimizing for rapid iteration

---

## Design Option 4: Semantic Phases (Action-Oriented)

**Philosophy**: Phase names describe the action being performed.

### Phases

#### Verify
Verify code quality without building.

```
mkdocs-ci:
  - Verify()
    - Prose linting
    - Formatting checks
    - Static analysis

miele-ci:
  - Verify()
    - Linting
    - Type checking
    - Formatting
```

#### Build
Build and test the artifact.

```
mkdocs-ci:
  - Build()
    - Build MkDocs site
    - Link validation
    - Asset optimization

miele-ci:
  - Build()
    - Vite build
    - Run tests
    - Bundle optimization
```

#### Certify
Security and compliance certification.

```
mkdocs-ci:
  - Certify()
    - Container scanning
    - License compliance
    - Vulnerability assessment

miele-ci:
  - Certify()
    - Trivy scan
    - Dependency audit
    - SBOM generation
```

#### Publish
Publish to registry.

```
mkdocs-ci:
  - Publish()
    - Push to GHCR with manifest

miele-ci:
  - Publish()
    - Push to GHCR with manifest
```

#### Deploy
Deploy to target environment.

```
mkdocs-ci:
  - Deploy(target)
    - Render / Fly.io / Cloud Run

miele-ci:
  - Deploy(target)
    - Fly.io
```

### Strengths
- ✅ Action verbs are clear and direct
- ✅ Each phase has single responsibility
- ✅ "Certify" emphasizes security importance
- ✅ Granular - easy to run individual phases

### Weaknesses
- ⚠️ More phases to orchestrate
- ⚠️ Build + Publish split may feel artificial
- ⚠️ "Certify" might sound too formal

### When to Use
- Organizations with security/compliance needs
- When you need fine-grained control
- Teams that prefer explicit over implicit
- Regulated industries

---

## Design Option 5: Dave Farley's Deployment Pipeline Pattern

**Philosophy**: Directly implements Farley's classic deployment pipeline stages.

### Phases

#### Commit Stage
Fast build and test - optimized for < 5 minutes.

```
mkdocs-ci:
  - CommitStage()
    - Fast linters only
    - Build site
    - Critical link checks only

miele-ci:
  - CommitStage()
    - Fast linters
    - Build app
    - Unit tests only
```

#### Acceptance Stage
Functional and integration tests.

```
mkdocs-ci:
  - AcceptanceStage()
    - Full link validation
    - Comprehensive linting
    - Visual regression (optional)

miele-ci:
  - AcceptanceStage()
    - Integration tests
    - E2E tests
    - API contract tests
```

#### Capacity Stage
Performance and load testing.

```
mkdocs-ci:
  - CapacityStage()
    - Load testing nginx
    - Performance budgets
    - Lighthouse CI

miele-ci:
  - CapacityStage()
    - Frontend performance
    - Bundle size checks
    - Lighthouse CI
```

#### Production Stage
Deploy to production with progressive rollout.

```
mkdocs-ci:
  - ProductionStage()
    - Blue/green deployment
    - Canary release
    - Monitoring validation

miele-ci:
  - ProductionStage()
    - Progressive Fly.io rollout
    - Health monitoring
```

### Strengths
- ✅ Faithful to CD literature
- ✅ Separates functional from non-functional testing
- ✅ Emphasizes production readiness
- ✅ Well-documented pattern

### Weaknesses
- ⚠️ "Capacity Stage" may be unclear
- ⚠️ Requires more infrastructure
- ⚠️ May be overkill for simple sites/apps
- ⚠️ Academic terminology

### When to Use
- Teams reading Farley's books
- Mature CD practices
- Complex applications needing performance validation
- When you want to follow established patterns

---

## Recommendations

### For mkdocs-ci
**Recommended: Design Option 2 (Validation Pipeline)**

Rationale:
- Documentation sites benefit from explicit validation phases
- Link checking is a distinct quality gate
- Security scanning containers is important
- Clear progression: code → build → security → release

### For miele-ci
**Recommended: Design Option 3 (Progressive Confidence)**

Rationale:
- Frontend development benefits from fast feedback
- Time budgets help developers understand wait times
- Less formal terminology fits modern dev culture
- Simple progression keeps pipeline maintainable

### General Guidance

**Use Design 1 (Classic)** if:
- Your team already uses "integration/delivery/deployment" terminology
- You have formal release processes
- You need clear artifact vs deployment separation

**Use Design 2 (Validation)** if:
- Security and compliance are priorities
- You want explicit quality gates
- Granular failure reporting is important

**Use Design 3 (Progressive Confidence)** if:
- Developer experience is paramount
- Fast feedback is critical
- You prefer simplicity

**Use Design 4 (Semantic)** if:
- You want fine-grained control
- Single responsibility per phase matters
- You're in a regulated industry

**Use Design 5 (Farley's Pattern)** if:
- You're implementing CD by the book
- You have capacity for comprehensive testing
- You need production deployment sophistication

---

## Terminology Clarity

### Do People Understand These Words?

| Term | Understanding | Notes |
|------|--------------|-------|
| **Integration** | ⚠️ Medium | Ambiguous - integrate with what? |
| **Delivery** | ⚠️ Medium | Often confused with deployment |
| **Deployment** | ✅ High | Well understood |
| **Validate** | ✅ High | Clear action, obvious purpose |
| **Build** | ✅ High | Universal understanding |
| **Test** | ✅ High | Clear and direct |
| **Publish** | ✅ High | Registry/artifact context clear |
| **Ship** | ✅ High | Modern, developer-friendly |
| **Certify** | ⚠️ Medium | May sound too formal |
| **Commit Stage** | ⚠️ Low | Requires CD knowledge |
| **Acceptance Stage** | ⚠️ Low | Academic terminology |
| **Capacity Stage** | ❌ Low | Unclear without explanation |
| **Quick/Full Check** | ✅ High | Intuitive time implication |

### Making It Practical

**Most Developer-Friendly:**
```
1. Check (quick)
2. Build
3. Scan (security)
4. Publish
5. Deploy
```

**Most Formal/Enterprise:**
```
1. Validate Code
2. Validate Build
3. Validate Security
4. Release
```

**Best Balance:**
```
1. Verify
2. Build
3. Test
4. Publish
5. Deploy
```

---

## Implementation Strategy

### Phase 1: Create Base Interfaces
Define common pipeline functions all modules should implement:

```go
// Every CI module should implement these phases
type Pipeline interface {
    Verify(ctx context.Context) error          // Fast quality checks
    Build(ctx context.Context) *Directory       // Build artifact
    Test(ctx context.Context) error            // Run tests
    Scan(ctx context.Context) error            // Security scan
    Publish(ctx context.Context) (string, error) // Push to registry
    Deploy(ctx context.Context, env string) error // Deploy to environment
}
```

### Phase 2: Refactor Existing Modules
Split monolithic functions into discrete phases:

```go
// Old: TestBuildPublishDeploy does everything
// New: Composable phases
func (m *MkdocsCi) Pipeline(ctx context.Context) error {
    if err := m.Verify(ctx); err != nil { return err }
    if err := m.Test(ctx); err != nil { return err }
    artifact := m.Build(ctx)
    if err := m.Scan(ctx, artifact); err != nil { return err }
    addr, err := m.Publish(ctx, artifact)
    if err != nil { return err }
    return m.Deploy(ctx, "production", addr)
}
```

### Phase 3: Create Orchestration Module
Build a `pipeline` module that composes these phases:

```bash
dagger init --sdk=go --name=pipeline pipeline
```

This module would orchestrate mkdocs-ci, miele-ci, and future modules with consistent phase execution.

---

## Conclusion

**Is it practical?** Yes, with caveats:

✅ **Practical aspects:**
- Provides clear structure
- Enables better error reporting
- Allows selective phase execution
- Improves caching and parallelization

⚠️ **Considerations:**
- Requires refactoring existing code
- Need to choose terminology carefully
- Must balance granularity vs complexity
- Documentation is critical for adoption

**Best approach**: Start with Design Option 3 (Progressive Confidence) for developer-facing modules, migrate to Design Option 2 (Validation Pipeline) as projects mature and require more formal quality gates.
