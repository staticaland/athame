// A Dagger module for Mermaid CLI (minlag/mermaid-cli)
//
// This module provides access to the Mermaid CLI for generating diagrams
// from Mermaid syntax.

package main

import (
	"fmt"

	"dagger/mermaid-cli/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=minlag/mermaid-cli
	// +default="11.12.0@sha256:bad64c9d9ad917c8dfbe9d9e9c162b96f6615ff019b37058638d16eb27ce7783"
	imageTag string,
) *MermaidCli {
	return &MermaidCli{
		ImageTag: imageTag,
	}
}

type MermaidCli struct {
	ImageTag string
}

// Base returns the base container with Mermaid CLI installed
func (m *MermaidCli) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("minlag/mermaid-cli:%s", m.ImageTag)).
		WithoutEntrypoint()
}
