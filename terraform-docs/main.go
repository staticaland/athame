// A Dagger module for terraform-docs

package main

import (
	"fmt"

	"dagger/terraform-docs/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=quay.io/terraform-docs/terraform-docs
	// +default="0.20.0@sha256:37329e2dc2518e7f719a986a3954b10771c3fe000f50f83fd4d98d489df2eae2"
	imageTag string,
) *TerraformDocs {
	return &TerraformDocs{
		ImageTag: imageTag,
	}
}

type TerraformDocs struct {
	ImageTag string
}

// Base returns the base container with terraform-docs installed
func (m *TerraformDocs) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("quay.io/terraform-docs/terraform-docs:%s", m.ImageTag)).
		WithoutEntrypoint()
}
