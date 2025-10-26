# Crane Image Finder

A Claude skill for finding and managing container images using the Crane CLI.

## Purpose

This skill helps you discover container image tags from container registries and ensures proper Renovate configuration for automated dependency updates.

## Usage

When asked to find container images or check available tags, use the `crane ls` command:

```bash
crane ls <registry>/<repository>/<image>
```

### Examples

List all tags for a container image:
```bash
crane ls gcr.io/projectcontour/contour
```

List with full references:
```bash
crane ls --full-ref docker.io/library/nginx
```

List repositories in a registry:
```bash
crane catalog <registry>
```

## Renovate Integration

**IMPORTANT**: When working with container images in any configuration files (Dockerfiles, YAML, etc.), always include Renovate regex patterns to enable automated dependency updates.

### Renovate Regex Pattern Template

Always use this pattern when adding or updating container images:

```dockerfile
# renovate: datasource=docker depName=<image-name>
ARG IMAGE_VERSION=<version>
FROM <image-name>:${IMAGE_VERSION}
```

### Common Renovate Patterns

#### 1. Dockerfile with ARG
```dockerfile
# renovate: datasource=docker depName=golang
ARG GO_VERSION=1.21.0
FROM golang:${GO_VERSION}
```

#### 2. Dockerfile with GitHub Releases
```dockerfile
# renovate: datasource=github-releases depName=aquasecurity/trivy
ARG TRIVY_VERSION=0.40.0
```

#### 3. YAML Configuration
```yaml
# renovate: datasource=docker depName=nginx versioning=docker
image: nginx:1.25.0
```

#### 4. With Digest Support
```dockerfile
# renovate: datasource=docker depName=alpine versioning=docker
FROM alpine:3.18@sha256:abc123...
```

### Renovate Configuration Example

For custom regex managers in `renovate.json`:

```json
{
  "regexManagers": [
    {
      "fileMatch": ["(^|/)Dockerfile$"],
      "matchStrings": [
        "# renovate: datasource=(?<datasource>.*?) depName=(?<depName>.*?)( versioning=(?<versioning>.*?))?\\s(?:ARG .*?_VERSION=|FROM .*?:)(?<currentValue>.*)\\s"
      ],
      "versioningTemplate": "{{#if versioning}}{{{versioning}}}{{else}}semver{{/if}}"
    }
  ]
}
```

## Key Datasources

- `datasource=docker` - For Docker Hub and container registries
- `datasource=github-releases` - For GitHub release versions
- `datasource=github-tags` - For GitHub tags
- `datasource=helm` - For Helm charts

## Versioning Templates

- `versioning=docker` - Docker-style versioning (default for container images)
- `versioning=semver` - Semantic versioning
- `versioning=pep440` - Python PEP 440 versioning
- `versioning=regex` - Custom regex-based versioning

## Best Practices

1. **Always include Renovate comments** when adding container images
2. **Use ARG for versions** to make updates easier
3. **Specify the datasource** explicitly for clarity
4. **Include versioning** when using non-standard version formats
5. **Test with crane** before adding to production configurations

## Crane Options

- `--full-ref` - Print full reference with registry and repository
- `--omit-digest-tags` - Don't list tags that represent digests
- `--platform` - Specify platform (e.g., linux/amd64)

## Resources

- [Crane Documentation](https://github.com/google/go-containerregistry/blob/main/cmd/crane/README.md)
- [Renovate Regex Manager](https://docs.renovatebot.com/modules/manager/regex/)
- [Renovate Docker Datasource](https://docs.renovatebot.com/modules/datasource/docker/)
