# Setting Up a Base Container for a Module

All container-based modules should provide a `Base()` function that returns a configured container.

## Purpose

- Provides a reusable container foundation
- Other module functions build on top of this base
- Centralizes cache volume configuration

## Example

```go
func (m *Terraform) Base() *dagger.Container {
    return dag.Container().
        From(fmt.Sprintf("hashicorp/terraform:%s", m.ImageTag)).
        WithoutEntrypoint()
}
```

## Usage in Other Functions

```go
func (m *Terraform) Plan(ctx context.Context, source *dagger.Directory) (string, error) {
    return m.Base().
        WithMountedDirectory("/src", source).
        WithWorkdir("/src").
        WithExec([]string{"terraform", "plan"}).
        Stdout(ctx)
}
```
