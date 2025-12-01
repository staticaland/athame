# Configuring Container Image Versions

All modules that wrap a container image MUST configure the image via an `imageTag` parameter in the constructor.

## Requirements

- **Always use digests** - never tags alone
- **Format**: `version@sha256:digest`
- **Add renovate comment** for automated updates
- **Store in struct** - add `ImageTag string` field

## Example

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

func (m *Terraform) Base() *dagger.Container {
    return dag.Container().
        From(fmt.Sprintf("hashicorp/terraform:%s", m.ImageTag)).
        WithoutEntrypoint()
}
```

## Finding Image Digests

Use the `container-image-lookup` subagent. Launch it FIRST so it runs in background while you write code.
