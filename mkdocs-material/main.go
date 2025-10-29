// A Dagger module for MkDocs Material
//
// This module provides functions for working with MkDocs Material,
// a Material Design theme for MkDocs static site generator.

package main

import (
	"dagger/mkdocs-material/internal/dagger"
	"fmt"
)

func New(
	// renovate: datasource=docker depName=squidfunk/mkdocs-material
	// +default="9.6.22@sha256:f5c556a6d30ce0c1c0df10e3c38c79bbcafdaea4b1c1be366809d0d4f6f9d57f"
	imageTag string,
) *MkdocsMaterial {
	return &MkdocsMaterial{
		ImageTag: imageTag,
	}
}

type MkdocsMaterial struct {
	ImageTag string
}

// Base returns the base container with MkDocs Material installed
func (m *MkdocsMaterial) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("squidfunk/mkdocs-material:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// Build builds the MkDocs documentation and returns the site directory
func (m *MkdocsMaterial) Build(
	// +defaultPath="/"
	source *dagger.Directory,
) *dagger.Directory {
	return m.Base().
		WithMountedDirectory("/docs", source).
		WithWorkdir("/docs").
		WithExec([]string{"build"}).
		Directory("/docs/site")
}
