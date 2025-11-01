# Athame Module & Pipeline Suggestions

This document contains ideas for new Dagger modules and demo pipelines, inspired by the existing codebase patterns.

## Overview of Existing Modules

The repository currently contains modules in these categories:

- **Infrastructure & IaC**: terraform, terraform-docs, localstack, aws-cli
- **Development**: node, uv, asdf
- **Documentation**: mkdocs-material, mermaid-cli
- **DevOps**: renovate, release-please, github-cli
- **Notifications**: apprise, ntfy
- **Utilities**: oras, httpie, boilerplate
- **Demos**: localstack-demo, hellos

## New Module Suggestions

### Security & Scanning

#### trivy
**Description**: Comprehensive vulnerability scanner for containers, filesystems, and IaC
**Use Cases**:
- Scan container images for CVEs
- Scan IaC configurations (Terraform, CloudFormation, Kubernetes)
- Scan filesystems and repositories
- Generate SBOM (Software Bill of Materials)

**Example Functions**:
```go
func (m *Trivy) ScanImage(imageRef string) *Container
func (m *Trivy) ScanFilesystem(source *Directory) *Container
func (m *Trivy) ScanIaC(source *Directory) *Container
func (m *Trivy) GenerateSBOM(imageRef string) *File
```

#### checkov
**Description**: Static code analysis for IaC security and compliance
**Use Cases**:
- Scan Terraform, CloudFormation, Kubernetes manifests
- Check against CIS benchmarks
- Custom policy enforcement
- CI/CD security gates

**Example Functions**:
```go
func (m *Checkov) Scan(source *Directory) *Container
func (m *Checkov) ScanWithFramework(source *Directory, framework string) *Container
```

#### gitleaks
**Description**: Secret scanning tool to prevent credentials from being committed
**Use Cases**:
- Pre-commit secret scanning
- Repository audit for exposed secrets
- CI/CD secret validation

**Example Functions**:
```go
func (m *Gitleaks) Detect(source *Directory) (string, error)
func (m *Gitleaks) Protect(source *Directory) (string, error)
```

#### cosign
**Description**: Container image signing and verification
**Use Cases**:
- Sign container images
- Verify image signatures
- Keyless signing with OIDC
- Attestation generation

**Example Functions**:
```go
func (m *Cosign) Sign(imageRef string, key *Secret) *Container
func (m *Cosign) Verify(imageRef string, key *Secret) (string, error)
func (m *Cosign) SignKeyless(imageRef string) *Container
```

#### syft
**Description**: SBOM generation tool for container images and filesystems
**Use Cases**:
- Generate CycloneDX or SPDX SBOMs
- Dependency tracking
- Compliance reporting

**Example Functions**:
```go
func (m *Syft) GenerateSBOM(imageRef string, format string) *File
func (m *Syft) ScanDirectory(source *Directory) *File
```

### Container & Registry Tools

#### crane
**Description**: Fast container image manipulation (already used in the repo!)
**Use Cases**:
- Copy images between registries
- Inspect image manifests
- Get image digests
- Export/import images

**Example Functions**:
```go
func (m *Crane) Copy(src string, dst string, credentials *Secret) *Container
func (m *Crane) Digest(imageRef string) (string, error)
func (m *Crane) Manifest(imageRef string) (string, error)
func (m *Crane) Export(imageRef string) *File
```

#### skopeo
**Description**: Container image operations without requiring root
**Use Cases**:
- Copy images between registries
- Inspect remote images
- Delete image tags
- Convert image formats

**Example Functions**:
```go
func (m *Skopeo) Copy(src string, dst string) *Container
func (m *Skopeo) Inspect(imageRef string) (string, error)
func (m *Skopeo) Delete(imageRef string) *Container
```

#### dive
**Description**: Explore and analyze container image layers
**Use Cases**:
- Image layer analysis
- Identify wasted space
- Optimize image size
- CI mode for image efficiency enforcement

**Example Functions**:
```go
func (m *Dive) Analyze(imageRef string) (string, error)
func (m *Dive) CI(imageRef string, highestUserWastedPercent float64) (string, error)
```

### Infrastructure & Cloud

#### pulumi
**Description**: Infrastructure as Code using programming languages
**Use Cases**:
- Multi-cloud infrastructure deployment
- State management
- Preview and update infrastructure

**Example Functions**:
```go
func (m *Pulumi) Preview(source *Directory, stack string) (string, error)
func (m *Pulumi) Up(source *Directory, stack string) (string, error)
func (m *Pulumi) Destroy(source *Directory, stack string) (string, error)
```

#### ansible
**Description**: Configuration management and automation
**Use Cases**:
- Server configuration
- Application deployment
- Multi-tier orchestration

**Example Functions**:
```go
func (m *Ansible) Playbook(source *Directory, playbook string, inventory *File) *Container
func (m *Ansible) Lint(source *Directory) (string, error)
```

#### kubectl
**Description**: Kubernetes command-line tool
**Use Cases**:
- Deploy applications to Kubernetes
- Manage cluster resources
- Debug running pods

**Example Functions**:
```go
func (m *Kubectl) Apply(manifests *Directory, kubeconfig *Secret) *Container
func (m *Kubectl) Get(resource string, kubeconfig *Secret) (string, error)
func (m *Kubectl) Rollout(deployment string, kubeconfig *Secret) *Container
```

#### helm
**Description**: Kubernetes package manager
**Use Cases**:
- Deploy applications using Helm charts
- Manage releases
- Template rendering

**Example Functions**:
```go
func (m *Helm) Install(chart string, release string, values *File, kubeconfig *Secret) *Container
func (m *Helm) Template(chart string, values *File) (string, error)
func (m *Helm) Upgrade(release string, chart string, kubeconfig *Secret) *Container
```

#### vault
**Description**: HashiCorp Vault for secrets management
**Use Cases**:
- Secrets storage and retrieval
- Dynamic credentials
- Encryption as a service

**Example Functions**:
```go
func (m *Vault) Read(path string, token *Secret) (string, error)
func (m *Vault) Write(path string, data string, token *Secret) *Container
```

### Development & Build Tools

#### go-sdk
**Description**: Go development environment and tools
**Use Cases**:
- Build Go applications
- Run tests with coverage
- Lint and format code

**Example Functions**:
```go
func (m *Go) Build(source *Directory, output string) *File
func (m *Go) Test(source *Directory) (string, error)
func (m *Go) ModTidy(source *Directory) *Directory
func (m *Go) Lint(source *Directory) (string, error)
```

#### python-sdk
**Description**: Python development environment
**Use Cases**:
- Build Python applications
- Run pytest
- Lint and format code

**Example Functions**:
```go
func (m *Python) Test(source *Directory) (string, error)
func (m *Python) Lint(source *Directory) (string, error)
func (m *Python) Build(source *Directory) *Directory
```

#### rust
**Description**: Rust development environment
**Use Cases**:
- Build Rust applications
- Run cargo tests
- Clippy linting

**Example Functions**:
```go
func (m *Rust) Build(source *Directory, release bool) *File
func (m *Rust) Test(source *Directory) (string, error)
func (m *Rust) Clippy(source *Directory) (string, error)
```

#### buf
**Description**: Protobuf tooling and linting
**Use Cases**:
- Lint protobuf files
- Generate code from proto files
- Breaking change detection

**Example Functions**:
```go
func (m *Buf) Lint(source *Directory) (string, error)
func (m *Buf) Generate(source *Directory) *Directory
func (m *Buf) Breaking(source *Directory, against string) (string, error)
```

#### just
**Description**: Command runner similar to make but simpler
**Use Cases**:
- Run project-specific commands
- Task automation
- Development workflows

**Example Functions**:
```go
func (m *Just) Run(source *Directory, recipe string) *Container
func (m *Just) List(source *Directory) (string, error)
```

### Documentation & Visualization

#### plantuml
**Description**: UML diagram generation from text
**Use Cases**:
- Generate sequence diagrams
- Create class diagrams
- Architecture visualization

**Example Functions**:
```go
func (m *PlantUML) Generate(source *File) *File
func (m *PlantUML) GenerateAll(source *Directory) *Directory
```

#### graphviz
**Description**: Graph visualization software
**Use Cases**:
- Dependency graphs
- Network diagrams
- Flow charts

**Example Functions**:
```go
func (m *Graphviz) Render(dotFile *File, format string) *File
```

#### hugo
**Description**: Fast static site generator
**Use Cases**:
- Build documentation sites
- Generate blogs
- Create landing pages

**Example Functions**:
```go
func (m *Hugo) Build(source *Directory) *Directory
func (m *Hugo) Serve(source *Directory, port int) *Service
```

#### asciidoctor
**Description**: AsciiDoc processor for technical documentation
**Use Cases**:
- Generate HTML from AsciiDoc
- Create PDFs
- Technical documentation

**Example Functions**:
```go
func (m *Asciidoctor) ConvertToHTML(source *Directory) *Directory
func (m *Asciidoctor) ConvertToPDF(source *Directory) *Directory
```

#### swagger-codegen / openapi-generator
**Description**: Generate API clients and documentation from OpenAPI specs
**Use Cases**:
- Generate API documentation
- Create client SDKs
- Server stubs

**Example Functions**:
```go
func (m *OpenAPIGenerator) GenerateDocs(spec *File) *Directory
func (m *OpenAPIGenerator) GenerateClient(spec *File, language string) *Directory
```

### Testing & Quality

#### shellcheck
**Description**: Shell script static analysis
**Use Cases**:
- Lint shell scripts
- Detect common errors
- Best practice enforcement

**Example Functions**:
```go
func (m *ShellCheck) Check(source *Directory) (string, error)
func (m *ShellCheck) CheckFile(script *File) (string, error)
```

#### hadolint
**Description**: Dockerfile linter
**Use Cases**:
- Dockerfile best practices
- Security checks
- Build optimization hints

**Example Functions**:
```go
func (m *Hadolint) Lint(dockerfile *File) (string, error)
```

#### yamllint
**Description**: YAML file linter
**Use Cases**:
- Validate YAML syntax
- Enforce YAML style
- Detect common errors

**Example Functions**:
```go
func (m *YAMLLint) Lint(source *Directory) (string, error)
```

#### golangci-lint
**Description**: Go linter aggregator
**Use Cases**:
- Run multiple Go linters
- Custom linter configuration
- Fast parallel execution

**Example Functions**:
```go
func (m *GolangCILint) Run(source *Directory) (string, error)
func (m *GolangCILint) RunWithConfig(source *Directory, config *File) (string, error)
```

#### prettier
**Description**: Opinionated code formatter
**Use Cases**:
- Format JavaScript/TypeScript
- Format JSON, YAML, Markdown
- Enforce consistent style

**Example Functions**:
```go
func (m *Prettier) Format(source *Directory) *Directory
func (m *Prettier) Check(source *Directory) (string, error)
```

#### actionlint
**Description**: GitHub Actions workflow linter
**Use Cases**:
- Validate workflow syntax
- Detect workflow errors
- Best practices

**Example Functions**:
```go
func (m *ActionLint) Lint(workflows *Directory) (string, error)
```

### Git & Version Control

#### git-cliff
**Description**: Changelog generator based on conventional commits
**Use Cases**:
- Generate changelogs
- Release notes
- Version bumping

**Example Functions**:
```go
func (m *GitCliff) Generate(source *Directory) *File
func (m *GitCliff) GenerateLatest(source *Directory) (string, error)
```

#### commitizen
**Description**: Conventional commits helper
**Use Cases**:
- Interactive commit message creation
- Enforce commit conventions
- Changelog generation

**Example Functions**:
```go
func (m *Commitizen) Bump(source *Directory) (string, error)
func (m *Commitizen) Changelog(source *Directory) *File
```

#### semantic-release
**Description**: Fully automated version management and package publishing
**Use Cases**:
- Automated semantic versioning
- Release automation
- Multi-channel releases

**Example Functions**:
```go
func (m *SemanticRelease) Release(source *Directory, token *Secret) *Container
func (m *SemanticRelease) DryRun(source *Directory) (string, error)
```

### Database

#### postgres
**Description**: PostgreSQL database service
**Use Cases**:
- Integration testing
- Development databases
- Schema migrations

**Example Functions**:
```go
func (m *Postgres) Run(version string) *Service
func (m *Postgres) WithInitScript(script *File) *Service
```

#### redis
**Description**: Redis in-memory data store
**Use Cases**:
- Caching layer for tests
- Session storage
- Queue backend

**Example Functions**:
```go
func (m *Redis) Run(version string) *Service
func (m *Redis) WithConfig(config *File) *Service
```

#### mongodb
**Description**: MongoDB NoSQL database
**Use Cases**:
- Document database for tests
- Development environment
- Data fixtures

**Example Functions**:
```go
func (m *MongoDB) Run(version string) *Service
func (m *MongoDB) WithInitData(data *Directory) *Service
```

### Communication & Collaboration

#### slack-cli
**Description**: Slack API interactions
**Use Cases**:
- Send build notifications
- Post to channels
- Upload files

**Example Functions**:
```go
func (m *SlackCLI) PostMessage(channel string, message string, token *Secret) *Container
func (m *SlackCLI) UploadFile(channel string, file *File, token *Secret) *Container
```

#### discord-webhook
**Description**: Discord webhook notifications
**Use Cases**:
- Build status notifications
- Deploy notifications
- Alert messages

**Example Functions**:
```go
func (m *Discord) Send(webhookURL *Secret, message string) *Container
func (m *Discord) SendEmbed(webhookURL *Secret, title string, description string, color string) *Container
```

### API & Processing

#### jq
**Description**: JSON processor
**Use Cases**:
- Parse JSON responses
- Transform JSON data
- Extract values from JSON

**Example Functions**:
```go
func (m *JQ) Process(input *File, filter string) (string, error)
func (m *JQ) ProcessString(input string, filter string) (string, error)
```

#### yq
**Description**: YAML/JSON/XML processor
**Use Cases**:
- Parse and transform YAML
- Convert between formats
- Extract configuration values

**Example Functions**:
```go
func (m *YQ) Process(input *File, expression string) (string, error)
func (m *YQ) Merge(files []*File) *File
```

#### grpcurl
**Description**: gRPC command-line client
**Use Cases**:
- Test gRPC services
- Inspect gRPC APIs
- Debug gRPC endpoints

**Example Functions**:
```go
func (m *GRPCurl) List(target string) (string, error)
func (m *GRPCurl) Call(target string, method string, data string) (string, error)
```

#### newman
**Description**: Postman collection runner
**Use Cases**:
- API testing in CI/CD
- Run Postman collections
- Generate test reports

**Example Functions**:
```go
func (m *Newman) Run(collection *File, environment *File) (string, error)
func (m *Newman) RunWithReporters(collection *File, reporters []string) *Directory
```

## Demo Pipeline Suggestions

### 1. Complete CI/CD Pipeline Demo

**Name**: `full-ci-demo`

**Description**: End-to-end CI/CD pipeline demonstrating multiple modules working together

**Pipeline Flow**:
```
1. Checkout → 2. Lint (shellcheck, hadolint, yamllint)
→ 3. Test (language-specific)
→ 4. Build (container image)
→ 5. Security Scan (trivy, checkov)
→ 6. Sign (cosign)
→ 7. SBOM (syft)
→ 8. Push to Registry
→ 9. Deploy
→ 10. Notify (ntfy/apprise)
```

**Example Function**:
```go
func (m *FullCIDemo) RunPipeline(
    ctx context.Context,
    source *Directory,
    registry string,
    credentials *Secret,
) (string, error)
```

### 2. Documentation Pipeline Demo

**Name**: `docs-pipeline-demo`

**Description**: Generate comprehensive documentation with diagrams

**Pipeline Flow**:
```
1. Generate Terraform Docs (terraform-docs)
→ 2. Generate API Docs (openapi-generator)
→ 3. Create Diagrams (mermaid-cli, plantuml)
→ 4. Build Site (mkdocs-material)
→ 5. Deploy to GitHub Pages
```

**Example Function**:
```go
func (m *DocsPipelineDemo) BuildDocs(
    ctx context.Context,
    source *Directory,
) *Directory
```

### 3. Security Scanning Pipeline Demo

**Name**: `security-scan-demo`

**Description**: Comprehensive security scanning for code and infrastructure

**Pipeline Flow**:
```
1. Secret Scanning (gitleaks)
→ 2. Code Scanning (specific linters)
→ 3. Container Scanning (trivy)
→ 4. IaC Scanning (checkov)
→ 5. Dependency Scanning (syft + vulnerability check)
→ 6. Generate Security Report
→ 7. Notify if issues found
```

**Example Function**:
```go
func (m *SecurityScanDemo) FullScan(
    ctx context.Context,
    source *Directory,
    imageRef string,
) (*Directory, error) // Returns report directory
```

### 4. Multi-Cloud Infrastructure Demo

**Name**: `multi-cloud-demo`

**Description**: Deploy infrastructure to multiple cloud providers using LocalStack

**Pipeline Flow**:
```
1. LocalStack (AWS services)
→ 2. Terraform Plan (multi-provider)
→ 3. Terraform Apply
→ 4. Run Tests against infrastructure
→ 5. Generate Docs
→ 6. Destroy
```

**Example Function**:
```go
func (m *MultiCloudDemo) DeployToLocalStack(
    ctx context.Context,
    terraformSource *Directory,
) (string, error)
```

### 5. Release Automation Demo

**Name**: `release-automation-demo`

**Description**: Automated release process with conventional commits

**Pipeline Flow**:
```
1. Analyze Commits (conventional commits)
→ 2. Generate Changelog (git-cliff)
→ 3. Bump Version (semantic-release)
→ 4. Create Release (release-please)
→ 5. Build & Tag Artifacts
→ 6. Push to Registry
→ 7. Create GitHub Release
→ 8. Notify (slack/discord)
```

**Example Function**:
```go
func (m *ReleaseAutomationDemo) Release(
    ctx context.Context,
    source *Directory,
    githubToken *Secret,
) (string, error)
```

### 6. Container Build & Sign Pipeline

**Name**: `container-secure-build-demo`

**Description**: Secure container build with signing and SBOM

**Pipeline Flow**:
```
1. Build Container
→ 2. Scan with Trivy
→ 3. Generate SBOM (syft)
→ 4. Sign Image (cosign)
→ 5. Push to Registry
→ 6. Verify Signature
→ 7. Generate Attestation
```

**Example Function**:
```go
func (m *ContainerSecureBuildDemo) BuildAndSign(
    ctx context.Context,
    dockerfile *File,
    context *Directory,
    imageRef string,
    signingKey *Secret,
) (string, error)
```

### 7. Monorepo Change Detection Demo

**Name**: `monorepo-demo`

**Description**: Detect changes in monorepo and build only affected modules

**Pipeline Flow**:
```
1. Detect Changed Paths
→ 2. Map Paths to Modules
→ 3. Build Affected Modules
→ 4. Test Affected Modules
→ 5. Deploy Changed Services
```

**Example Function**:
```go
func (m *MonorepoDemo) BuildAffected(
    ctx context.Context,
    source *Directory,
    baseBranch string,
) (string, error)
```

### 8. API Development Pipeline

**Name**: `api-dev-demo`

**Description**: Complete API development workflow

**Pipeline Flow**:
```
1. Validate OpenAPI Spec
→ 2. Generate Server Stubs
→ 3. Generate Client SDKs
→ 4. Run Integration Tests (newman/httpie)
→ 5. Generate API Documentation
→ 6. Deploy to LocalStack API Gateway
```

**Example Function**:
```go
func (m *APIDevDemo) DevelopAPI(
    ctx context.Context,
    openapiSpec *File,
) (*Directory, error)
```

### 9. Database Migration Demo

**Name**: `db-migration-demo`

**Description**: Database migration testing and validation

**Pipeline Flow**:
```
1. Start Database Service (postgres/mysql)
→ 2. Run Migrations
→ 3. Seed Test Data
→ 4. Run Integration Tests
→ 5. Generate Schema Documentation
→ 6. Cleanup
```

**Example Function**:
```go
func (m *DBMigrationDemo) TestMigrations(
    ctx context.Context,
    migrations *Directory,
) (string, error)
```

### 10. IaC Testing Pipeline

**Name**: `iac-testing-demo`

**Description**: Comprehensive infrastructure testing

**Pipeline Flow**:
```
1. Terraform Format Check
→ 2. Terraform Validate
→ 3. Security Scan (checkov)
→ 4. Plan (terraform-docs output)
→ 5. Apply to LocalStack
→ 6. Run Compliance Tests
→ 7. Generate Documentation
→ 8. Destroy
```

**Example Function**:
```go
func (m *IaCTestingDemo) TestInfrastructure(
    ctx context.Context,
    terraformSource *Directory,
) (*Directory, error) // Returns test report + docs
```

### 11. Microservices Integration Demo

**Name**: `microservices-integration-demo`

**Description**: Test multiple services working together

**Pipeline Flow**:
```
1. Start Dependencies (postgres, redis, localstack)
→ 2. Build Service Containers
→ 3. Start Services with Service Bindings
→ 4. Run Integration Tests
→ 5. Collect Logs
→ 6. Generate Coverage Report
```

**Example Function**:
```go
func (m *MicroservicesDemo) IntegrationTest(
    ctx context.Context,
    services map[string]*Directory,
) (string, error)
```

### 12. Compliance & Audit Demo

**Name**: `compliance-audit-demo`

**Description**: Generate compliance reports and audit trails

**Pipeline Flow**:
```
1. Scan Code (multiple linters)
→ 2. Check Licenses (SBOM analysis)
→ 3. Security Scan (trivy, checkov)
→ 4. Generate Compliance Report
→ 5. Archive Reports
→ 6. Notify Stakeholders
```

**Example Function**:
```go
func (m *ComplianceAuditDemo) GenerateReport(
    ctx context.Context,
    source *Directory,
    images []string,
) (*Directory, error)
```

## Implementation Priority

### High Priority (Most Valuable)

1. **trivy** - Essential for security scanning
2. **crane** - Already used in tooling, should be a module
3. **go-sdk** - Common language in Dagger ecosystem
4. **kubectl** - Kubernetes is widely used
5. **jq/yq** - Data processing essentials
6. **full-ci-demo** - Showcases module composition

### Medium Priority

1. **checkov** - IaC security
2. **cosign** - Image signing
3. **syft** - SBOM generation
4. **helm** - Kubernetes package management
5. **postgres/redis** - Testing databases
6. **security-scan-demo** - Security showcase

### Low Priority (Nice to Have)

1. **plantuml** - Diagram generation
2. **hugo** - Static site alternative
3. **slack-cli** - Additional notification channel
4. **newman** - API testing
5. **compliance-audit-demo** - Specialized use case

## Module Development Patterns

Based on existing modules, follow these patterns:

1. **Container-based modules**: Use `imageTag` parameter in constructor with renovate comment
2. **Tool-based modules**: Use `asdf` for version management (like github-cli, boilerplate)
3. **Base() function**: Always provide a base container
4. **Service modules**: Return `*Service` for long-running processes (like localstack)
5. **Demo modules**: Show real-world integration patterns

## Next Steps

1. Review and prioritize module suggestions
2. Identify quick wins (modules with existing containers)
3. Plan demo pipelines that showcase multiple modules
4. Consider community feedback and use cases
5. Create issues for tracking implementation
