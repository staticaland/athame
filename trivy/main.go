// Wrapper for Trivy CLI
// Scans container images for vulnerabilities
// Uses official Trivy image

package main

import (
	"context"
	"fmt"
	"strconv"

	"dagger/trivy/internal/dagger"
)

func New(
	// renovate: datasource=docker depName=aquasec/trivy
	// +default="0.67.2@sha256:e2b22eac59c02003d8749f5b8d9bd073b62e30fefaef5b7c8371204e0a4b0c08"
	imageTag string,
) *Trivy {
	return &Trivy{
		ImageTag: imageTag,
	}
}

// Wrapper for Trivy CLI
type Trivy struct {
	ImageTag string
}

// Base returns the base container with Trivy installed
func (m *Trivy) Base() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("aquasec/trivy:%s", m.ImageTag)).
		WithMountedCache("/root/.cache/trivy", dag.CacheVolume("trivy-db-cache"))
}

// ScanImage scans a container image reference for vulnerabilities
func (m *Trivy) ScanImage(
	ctx context.Context,
	imageRef string,
	// +optional
	// +default="UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL"
	severity string,
	// +optional
	// +default=0
	exitCode int,
	// +optional
	// +default="table"
	format string,
) (string, error) {
	return m.Base().
		WithExec([]string{
			"trivy",
			"image",
			"--quiet",
			"--severity", severity,
			"--exit-code", strconv.Itoa(exitCode),
			"--format", format,
			imageRef,
		}).
		Stdout(ctx)
}

// ScanContainer scans a Dagger Container for vulnerabilities
func (m *Trivy) ScanContainer(
	ctx context.Context,
	ctr *dagger.Container,
	imageRef string,
	// +optional
	// +default="UNKNOWN,LOW,MEDIUM,HIGH,CRITICAL"
	severity string,
	// +optional
	// +default=0
	exitCode int,
	// +optional
	// +default="table"
	format string,
) (string, error) {
	return m.Base().
		WithMountedFile("/scan/"+imageRef, ctr.AsTarball()).
		WithExec([]string{
			"trivy",
			"image",
			"--quiet",
			"--severity", severity,
			"--exit-code", strconv.Itoa(exitCode),
			"--format", format,
			"--input", "/scan/" + imageRef,
		}).
		Stdout(ctx)
}
