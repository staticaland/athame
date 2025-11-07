# Dagger Development Guidelines

## Monorepo Structure

This is a Dagger monorepo where we create and manage multiple Dagger modules, always using the Go SDK.

### Creating New Modules

To create a new module, use:

```bash
dagger init --sdk=go --name=<module-name> <module-name>
```

Example:

```bash
dagger init --sdk=go --name=basics basics
```

Command format: `dagger init [options] [path]`

### Using Modules

To use a module in your project, use the `dagger install` command:

```bash
dagger install [options] <module>
```

### Exploring Modules

Before using a module, explore its available functions to understand its API:

```bash
dagger functions --mod <module-path>
```

Example:

```bash
dagger functions --mod ./terraform-docs
```

This shows all available functions, their parameters, and return types. Always check what functions exist before implementing - don't assume method names or signatures.

### Exploring Container Commands

To explore commands available in a container (e.g., checking CLI help text), use `with-exec` with `--args`:

```bash
# Explore command help
dagger call --mod ./terraform-docs base with-exec --args="terraform-docs,--help" stdout

# Check command options
dagger call --mod ./terraform-docs base with-exec --args="terraform-docs,markdown,--help" stdout
```

**Never use `terminal`** - it requires interactive TTY which is not supported in non-interactive contexts. Always use `with-exec` with comma-separated `--args` instead.

### Exporting Results

To export files or directories from Dagger functions to your local filesystem, use the `export` command:

```bash
dagger call <function-name> export --path=<local-path>
```

Example with the `terraform-docs` function in the repo module:

```bash
dagger call terraform-docs export --path=fixtures/terraform
```

### Directory Structure

- **`.dagger/`** - The repo orchestration module. This is the module that can be called from the repository root using `dagger call <function-name>`. Use `dagger install` here to compose multiple modules, or create dedicated wrapper modules elsewhere. Project CI may also live here.
- **Individual module directories** - Each module lives in its own directory at the repo root.

### Repo Module

The repo module is located in `.dagger/` and can be called directly from the repository root without the `--mod` flag:

```bash
# List repo module functions
dagger functions

# Call a repo module function
dagger call terraform-docs export --path ./fixtures/terraform
```

This module typically orchestrates other modules and provides high-level functions for the entire repository.

## Working Directory

**Always work from the repository root.** All Dagger commands should be run from the top-level directory of the repository.

### Path Patterns

Use these path patterns consistently:

- **`.` (single dot)** - represents the current directory (the repo root when you're working from there)
- **`./module-name`** - path to a specific module directory from the repo root
- **`./.dagger`** - path to the repo orchestration module

### When to Use Single Dot (`.`)

Use `.` in these contexts:

1. **When installing a module from the repo root into `.dagger/`:**

   ```bash
   dagger install --mod ./.dagger ./terraform-docs
   ```

   This installs the `terraform-docs` module into the repo module.

2. **When a module depends on the current directory context:**

   ```bash
   dagger call --mod ./terraform-docs generate --source=.
   ```

   This passes the repo root as a directory parameter.

3. **When working within a module that needs to reference its own directory** (rare, since we avoid `cd`).

### Never Use

- **`..` (parent directory)** - avoid relative parent paths; use explicit paths from repo root instead
- **`cd` commands** - stay in the repo root and use `--mod` with explicit paths

## Calling Dagger Modules

When calling Dagger modules, always follow these rules:

- **Stay at repository root** - run all commands from the top-level directory
- **Always use `--mod` flag** - explicitly specify the module path with `--mod` instead of relying on directory context
- **Never use `cd`** - avoid changing directories to call a module; use `--mod` with the full path instead
- **Exception: Repo module** - the repo module (`.dagger/`) can be called from the repository root without `--mod`

### Examples

```bash
# Good: explicit module path from repo root
dagger call --mod ./terraform-docs generate
dagger functions --mod ./terraform-docs
dagger install --mod ./.dagger ./terraform-docs

# Good: using single dot for directory parameters
dagger call --mod ./terraform-docs generate --source=.
dagger call --mod ./mkdocs-ci build --docs=./docs

# Good: repo module from repository root (no --mod needed)
dagger call terraform-docs
dagger functions

# Good: installing dependencies for a specific module
dagger install --mod ./terraform-docs github.com/example/some-dependency

# Bad: changing directory
cd terraform-docs && dagger call generate
cd terraform-docs && dagger functions

# Bad: omitting --mod for non-repo modules
dagger call generate
dagger functions

# Bad: using parent directory references
dagger install --mod ./.dagger ../terraform-docs  # Use ./terraform-docs instead

# Bad: using cd to get context
cd terraform-docs && dagger call generate --source=.  # Use --mod and --source=./terraform-docs
```

### Real-World Workflow

Here's a typical workflow from the repository root:

```bash
# 1. Create a new module (stay at repo root)
dagger init --sdk=go --name=my-module my-module

# 2. Explore what it provides (stay at repo root)
dagger functions --mod ./my-module

# 3. Install it into the repo module (stay at repo root)
dagger install --mod ./.dagger ./my-module

# 4. Call it with the repo root as source (stay at repo root)
dagger call --mod ./my-module process --source=.

# 5. Export results (stay at repo root)
dagger call --mod ./my-module process --source=. export --path=./output
```

Notice: no `cd` commands, all paths explicit from repo root.

## Import Path

Always import Dagger from:

```go
import "dagger/<module-name>/internal/dagger"
```

Replace `<module-name>` with your actual module name.

## Generated Files - Do Not Edit

These files are auto-generated by Dagger and should never be manually edited or read:

- `dagger.gen.go`
- `internal/` directory

## Constructors

Use the `func New()` pattern when multiple functions need the same parameters.

### When to Use

- **Shared configuration** - parameters used across multiple functions
- **Module-wide state** - avoid repeating the same parameters in every function

### Requirements

- **Public fields only** - capitalize field names in Go (e.g., `Greeting`, `Name`) for proper serialization
- **One constructor per module** - Dagger modules have only one constructor
- **Return module pointer** - return `*YourModule` from `New()`

### Container Image Configuration

**All modules that wrap a container image MUST configure the image reference via an `imageTag` parameter in the constructor.**

This is the standard pattern across all modules in this repository.

### Pattern

- **Always use imageTag parameter** - every container-based module constructor must have this parameter
- **Store version+digest** - use `version@sha256:digest` format (never use tags alone)
- **Add renovate comment** - place above the parameter for automated dependency updates
- **Store in struct** - add `ImageTag string` field to your module struct
- **Use in Base()** - construct full image reference with `fmt.Sprintf("image:%s", m.ImageTag)`

### Example

```go
func New(
    // renovate: datasource=docker depName=hashicorp/terraform
    // +default="1.13.4@sha256:eebc943e69008b6d6d986800087164274d8c92d83db8d53fb9baa4ccff309884"
    imageTag string,
) *Terraform {
    return &Terraform{
        ImageTag: imageTag,
    }
}

type Terraform struct {
    ImageTag string
}

// Base returns the base container with Terraform installed
func (m *Terraform) Base() *dagger.Container {
    return dag.Container().
        From(fmt.Sprintf("hashicorp/terraform:%s", m.ImageTag)).
        WithoutEntrypoint()
}
```

### General Constructor Example

```go
func New(
    // +default="https://api.example.com"
    apiUrl string,
    // +optional
    apiKey string,
) *MyModule {
    return &MyModule{
        ApiUrl: apiUrl,
        ApiKey: apiKey,
    }
}

type MyModule struct {
    ApiUrl string
    ApiKey string
}

func (m *MyModule) GetUser(id string) (*User, error) {
    // Use m.ApiUrl and m.ApiKey without passing them as parameters
}

func (m *MyModule) CreateUser(name string) (*User, error) {
    // Same fields available here
}
```

## Base Function

All modules should provide a `Base()` function that returns a configured base container.

### Purpose

- **Container foundation** - provides a reusable base container with runtime image and dependencies
- **Shared setup** - other module functions build on top of this base
- **Cache optimization** - centralizes cache volume configuration

### Requirements

- **Always use image digests** - never use tags alone (e.g., `node:lts-alpine`)
- **Prefer specific version tags** - use fully qualified versions including base OS version when available (e.g., `20.11.1-alpine3.19` over `20.11.1-alpine`).
- **Use crane for digests** - run `crane ls` and `crane digest` with the `Task` tool to get the latest image and digest
- **Add renovate comments** - enable automated dependency updates

### Example

```go
// Base returns the base container with runtime and dependencies installed
func (m *MyModule) Base() *dagger.Container {
    return dag.Container().
        // renovate: datasource=docker depName=node
        From("node:20.11.1-alpine3.19@sha256:...").
        WithoutEntrypoint().
        WithMountedDirectory("/src", m.Src).
        WithWorkdir("/src")
}
```

### Finding Image Digests

Use the `container-image-lookup` subagent to find the latest version and image digest.

**IMPORTANT: When it's obvious a module needs a container image digest/tag, launch the container-image-lookup subagent as the FIRST task.** This allows the lookup to happen in the background while you work on other parts of the module.

**Workflow:**

1. Launch the subagent using the Task tool to find the image tag and digest for the container image
2. Continue editing the module code without the tag and digest (use a placeholder or leave it empty)
3. When the subagent returns with the result (in `version@sha256:digest` format), add it to the constructor

This allows you to work on other parts of the module while the image lookup happens concurrently.

## Function Parameters

Use Dagger annotations for parameter handling:

- `+default="value"` - for parameters with default values
- `+optional` - for parameters that can be omitted

### Required Patterns

- **Use annotations exclusively** - let Dagger handle defaults, not your function body
- **Trust the framework** - write functions assuming defaults are already applied
- **Document well** - add descriptive comments for optional parameters

### Forbidden Patterns

- **No manual checks** - don't use `if param == ""` for parameters with `+default` annotations
- **No defensive logic** - avoid conditional logic for optional parameters unless your business logic requires it
- **No mixing approaches** - don't combine annotations with manual default handling

### Directory Parameters

For functions that work with repository source code, use `+defaultPath` to automatically resolve directory paths.

#### Path Resolution in Git Repositories

- **Absolute paths** (`/`, `/src`) - resolve from repository root
- **Relative path** (`.`) - resolves to the directory containing `dagger.json`
- **Parent path** (`..`) - resolves to parent of the directory containing `dagger.json`

#### Pattern

Use `+defaultPath="/"` for modules that should work with the entire repository by default:

```go
func (m *MyModule) ReadDir(
    ctx context.Context,
    // +defaultPath="/"
    source *dagger.Directory,
) ([]string, error) {
    return source.Entries(ctx)
}
```

## Example

```go
// Deploy runs the deployment process
// +optional
func (m *MyModule) Deploy(
    // +default="production"
    environment string,
    // +optional
    dryRun bool,
) error {
    // Function body assumes defaults are already applied
    // No need for: if environment == "" { environment = "production" }
}
```

## Useful Links

- [Dagger Go SDK API Reference](https://pkg.go.dev/dagger.io/dagger)
- [Dagger Documentation](https://docs.dagger.io/llms.txt)
