# Making Parameters Optional or Giving Them Default Values

Use Dagger annotations to handle parameter defaults and optionality.

## Annotations

- `+default="value"` - parameter has a default value
- `+optional` - parameter can be omitted
- `+defaultPath="/"` - directory parameter defaults to repo root

## Example

```go
func (m *MyModule) Deploy(
    // +default="production"
    environment string,
    // +optional
    dryRun bool,
) error {
    // Assume defaults already applied by Dagger
}
```

## Directory Parameters

```go
func (m *MyModule) Build(
    // +defaultPath="/"
    source *dagger.Directory,
) *dagger.Container {
    // source defaults to repository root
}
```

## Important

**Do not manually check for empty values** when using `+default`. Dagger applies defaults before your function runs.

Wrong:

```go
if environment == "" {
    environment = "production"  // Don't do this
}
```

Right:

```go
// Just use environment directly - Dagger already applied the default
```
