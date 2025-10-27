# Terraform Dagger Module

A Dagger module for Terraform operations.

## Usage

### Default Version

By default, the module uses the latest stable Terraform version:

```bash
dagger call base with-exec --args=terraform,version stdout
```

### Override Terraform Version

You can override the Terraform version using the `--image-tag` parameter:

```bash
dagger call --image-tag=1.12.0 base with-exec --args=terraform,version stdout
```

The `imageTag` parameter accepts any valid tag from the [hashicorp/terraform](https://hub.docker.com/r/hashicorp/terraform/tags) Docker image, optionally with a digest for pinning:

```bash
dagger call --image-tag=1.12.0@sha256:... base
```

## Functions

### Base()

Returns the base container with Terraform installed and ready to use.
