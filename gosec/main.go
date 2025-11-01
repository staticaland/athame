// A Dagger module for securego/gosec

package main

import (
	"fmt"

	"dagger/gosec/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=securego/gosec
	// +default="2.22.10@sha256:c8852d609f9af551387555a81808a3bca8d172629b124fab0d83c937cabc2f3d"
	imageTag string,
) *Gosec {
	return &Gosec{
		ImageTag: imageTag,
	}
}

type Gosec struct {
	ImageTag string
}

// Base returns the base container with gosec installed
func (m *Gosec) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("securego/gosec:%s", m.ImageTag)).
		WithoutEntrypoint()
}
