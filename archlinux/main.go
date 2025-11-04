// A Dagger module for Arch Linux containers

package main

import (
	"fmt"

	"dagger/archlinux/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=archlinux/archlinux
	// +default="base-20251019.0.437072@sha256:4524236733437ff1f35531147aa444b32f674d9f328aebe06d3511be575c80a3"
	imageTag string,
) *Archlinux {
	return &Archlinux{
		ImageTag: imageTag,
	}
}

type Archlinux struct {
	ImageTag string
}

// Base returns a base Arch Linux container
func (m *Archlinux) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("archlinux/archlinux:%s", m.ImageTag)).
		WithoutEntrypoint()
}

// WithPackages returns a container with the specified packages installed
func (m *Archlinux) WithPackages(packages []string) *dagger.Container {
	args := append([]string{"pacman", "-Sy", "--noconfirm"}, packages...)
	return m.Base().WithExec(args)
}
