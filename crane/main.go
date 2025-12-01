// A Dagger module for interacting with container registries using Crane
//
// Crane is a tool from google/go-containerregistry for interacting with
// remote images and registries. This module provides a Dagger interface
// to common Crane operations like listing tags, getting digests, copying
// images, and inspecting manifests.

package main

import (
	"context"
	"fmt"

	"dagger/crane/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=gcr.io/go-containerregistry/crane
	// +default="v0.20.3@sha256:fc86bcad43a000c2a1ca926a1e167db26c053cebc3fa5d14285c72773fb8c11d"
	imageTag string,
) *Crane {
	return &Crane{
		ImageTag: imageTag,
	}
}

type Crane struct {
	ImageTag string
}

// Base returns a configured Crane container
func (m *Crane) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("gcr.io/go-containerregistry/crane:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// List returns all tags for a given image repository
//
// Example: crane list gcr.io/go-containerregistry/crane
func (m *Crane) List(
	ctx context.Context,
	// Repository to list tags from (e.g., "gcr.io/go-containerregistry/crane")
	repository string,
) (string, error) {
	return m.Base().
		WithExec([]string{"crane", "ls", repository}).
		Stdout(ctx)
}

// Digest returns the digest of a specific image
//
// Example: crane digest gcr.io/go-containerregistry/crane:latest
func (m *Crane) Digest(
	ctx context.Context,
	// Full image reference (e.g., "gcr.io/go-containerregistry/crane:latest")
	image string,
) (string, error) {
	return m.Base().
		WithExec([]string{"crane", "digest", image}).
		Stdout(ctx)
}

// Manifest returns the manifest of a specific image
//
// Example: crane manifest gcr.io/go-containerregistry/crane:latest
func (m *Crane) Manifest(
	ctx context.Context,
	// Full image reference (e.g., "gcr.io/go-containerregistry/crane:latest")
	image string,
) (string, error) {
	return m.Base().
		WithExec([]string{"crane", "manifest", image}).
		Stdout(ctx)
}

// Config returns the config file of a specific image
//
// Example: crane config gcr.io/go-containerregistry/crane:latest
func (m *Crane) Config(
	ctx context.Context,
	// Full image reference (e.g., "gcr.io/go-containerregistry/crane:latest")
	image string,
) (string, error) {
	return m.Base().
		WithExec([]string{"crane", "config", image}).
		Stdout(ctx)
}

// Validate validates that an image reference exists
//
// Example: crane validate --remote gcr.io/go-containerregistry/crane:latest
func (m *Crane) Validate(
	ctx context.Context,
	// Full image reference (e.g., "gcr.io/go-containerregistry/crane:latest")
	image string,
) (string, error) {
	return m.Base().
		WithExec([]string{"crane", "validate", "--remote", image}).
		Stdout(ctx)
}

// Copy copies an image from one repository to another
//
// Example: crane copy source-image:tag dest-image:tag
func (m *Crane) Copy(
	ctx context.Context,
	// Source image reference (e.g., "gcr.io/source/image:tag")
	source string,
	// Destination image reference (e.g., "gcr.io/dest/image:tag")
	destination string,
	// +optional
	// Registry secret for authentication
	secret *dagger.Secret,
) (string, error) {
	ctr := m.Base()

	if secret != nil {
		ctr = ctr.WithMountedSecret("/root/.docker/config.json", secret)
	}

	return ctr.
		WithExec([]string{"crane", "copy", source, destination}).
		Stdout(ctx)
}

// Export exports image contents as a tarball
//
// Example: crane export image:tag output.tar
func (m *Crane) Export(
	ctx context.Context,
	// Full image reference (e.g., "gcr.io/go-containerregistry/crane:latest")
	image string,
) *dagger.File {
	return m.Base().
		WithExec([]string{"crane", "export", image, "/tmp/image.tar"}).
		File("/tmp/image.tar")
}

// Tag adds a tag to an existing image
//
// Example: crane tag image:old-tag new-tag
func (m *Crane) Tag(
	ctx context.Context,
	// Full image reference (e.g., "gcr.io/go-containerregistry/crane:latest")
	image string,
	// New tag to add
	tag string,
	// +optional
	// Registry secret for authentication
	secret *dagger.Secret,
) (string, error) {
	ctr := m.Base()

	if secret != nil {
		ctr = ctr.WithMountedSecret("/root/.docker/config.json", secret)
	}

	return ctr.
		WithExec([]string{"crane", "tag", image, tag}).
		Stdout(ctx)
}
