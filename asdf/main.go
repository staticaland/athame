// A Dagger module for working with asdf version manager
//
// This module provides functionality for managing runtime versions using asdf.

package main

import (
	"fmt"

	"dagger/asdf/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=asdfvm/asdf
	// +default="alpine-v0.17.0@sha256:9744fdf066a668d477186560e2680f87bc935d6f1f17d020c00db83e1006d187"
	imageTag string,
) *Asdf {
	return &Asdf{
		ImageTag: imageTag,
	}
}

type Asdf struct {
	ImageTag string
}

// Base returns the base container with asdf installed
func (m *Asdf) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("asdfvm/asdf:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// InstallPlugin installs an asdf plugin, installs the specified version, and sets it globally
func (m *Asdf) InstallPlugin(
	// Plugin name (e.g., "nodejs", "python", "github-cli")
	pluginName string,
	// Plugin repository URL
	pluginUrl string,
	// Version to install and set globally
	// +default="latest"
	version string,
) *dagger.Container {
	return m.Base().
		WithExec([]string{"asdf", "plugin", "add", pluginName, pluginUrl}).
		WithExec([]string{"asdf", "install", pluginName, version}).
		WithExec([]string{"asdf", "set", pluginName, version})
}
