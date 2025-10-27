// A Dagger module for working with Boilerplate
//
// This module provides functionality for using the Boilerplate templating tool
// via asdf version manager.

package main

import (
	"dagger/boilerplate/internal/dagger"
)

func New(
	// renovate: datasource=github-releases depName=gruntwork-io/boilerplate
	// +default="0.10.1"
	version string,
) *Boilerplate {
	return &Boilerplate{
		Version: version,
	}
}

type Boilerplate struct {
	Version string
}

// Base returns the base container with boilerplate installed via asdf
func (m *Boilerplate) Base() *dagger.Container {
	return dag.Asdf().Base().
		WithUser("root").
		WithExec([]string{"apk", "add", "--no-cache", "git", "openssh-client"}).
		WithUser("asdf").
		WithExec([]string{"asdf", "plugin", "add", "boilerplate", "https://github.com/gruntwork-io/asdf-boilerplate.git"}).
		WithExec([]string{"asdf", "install", "boilerplate", m.Version}).
		WithExec([]string{"asdf", "set", "boilerplate", m.Version})
}
