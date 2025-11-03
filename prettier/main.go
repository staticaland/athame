// A Dagger module for Prettier code formatting

package main

import (
	"fmt"

	"dagger/prettier/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=node
	// +default="22.21.1-alpine3.22@sha256:b2358485e3e33bc3a33114d2b1bdb18cdbe4df01bd2b257198eb51beb1f026c5"
	imageTag string,
) *Prettier {
	return &Prettier{
		ImageTag: imageTag,
	}
}

type Prettier struct {
	ImageTag string
}

// Base returns the base container with Node.js and Prettier installed globally
func (m *Prettier) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("node:%s", m.ImageTag)).
		WithoutEntrypoint().
		WithExec([]string{"npm", "install", "-g", "prettier"})
}
