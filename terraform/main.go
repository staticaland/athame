// A Dagger module for Terraform operations

package main

import (
	"fmt"

	"dagger/terraform/internal/dagger"
)

type Terraform struct {
	ImageTag string
}

func New(
	// renovate: datasource=docker depName=hashicorp/terraform
	// +default="1.13.4@sha256:eebc943e69008b6d6d986800087164274d8c92d83db8d53fb9baa4ccff309884"
	imageTag string,
) *Terraform {
	return &Terraform{
		ImageTag: imageTag,
	}
}

// Base returns the base container with Terraform installed
func (m *Terraform) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("hashicorp/terraform:%s", m.ImageTag)).
		WithoutEntrypoint()
}
