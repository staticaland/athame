// A Dagger module for HTTPie, a user-friendly HTTP client
//
// HTTPie is a command-line HTTP client with an intuitive UI, JSON support,
// syntax highlighting, and more.

package main

import (
	"dagger/httpie/internal/dagger"
	"fmt"
)

func New(
	// renovate: datasource=docker depName=alpine/httpie
	// +default="3.2.4@sha256:cd81ee5ddd4970cc3175fddf1fdfad8df909a473eb5f82547e37ab510ed62fc5"
	imageTag string,
) *Httpie {
	return &Httpie{
		ImageTag: imageTag,
	}
}

type Httpie struct {
	ImageTag string
}

// Base returns the base container with HTTPie installed
func (m *Httpie) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("alpine/httpie:%s", m.ImageTag)).
		WithoutEntrypoint()
}
