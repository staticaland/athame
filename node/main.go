// A Dagger module for Node.js development and build tasks

package main

import (
	"fmt"

	"dagger/node/internal/dagger"
)

type Node struct {
	// Source directory containing Node.js project files
	Src *dagger.Directory

	// Node.js image tag with version and digest
	ImageTag string
}

// New creates a new Node module instance
func New(
	// Source directory containing Node.js project files
	// +defaultPath="/"
	src *dagger.Directory,
	// renovate: datasource=docker depName=node
	// +default="22.21.0-alpine3.22@sha256:bd26af08779f746650d95a2e4d653b0fd3c8030c44284b6b98d701c9b5eb66b9"
	imageTag string,
) *Node {
	return &Node{
		Src:      src,
		ImageTag: imageTag,
	}
}

// Base returns the base container with Node.js runtime and dependencies installed
func (m *Node) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("node:%s", m.ImageTag)).
		WithoutEntrypoint().
		WithMountedDirectory("/src", m.Src).
		WithWorkdir("/src")
}
