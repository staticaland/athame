// Deploy containers to fly.io
//
// A utility module to deploy a container to fly.io using a very basic default configuration

package main

import (
	"context"
	"fmt"

	"dagger/flyio/internal/dagger"
)

type Flyio struct{}

// Deploy deploys an application to fly.io with an image ref
func (m *Flyio) Deploy(
	ctx context.Context,
	// The fly.io app name
	app string,
	// The container image reference to deploy
	image string,
	// The fly.io API token
	token *dagger.Secret,
	// Primary region for the app (see https://fly.io/docs/reference/regions/)
	// +default="arn"
	primaryRegion string,
	// +default=8080
	internalPort int,
) (string, error) {
	config := fmt.Sprintf(`app = "%s"
primary_region = "%s"

[build]
  image = "%s"

[http_service]
  internal_port = %d
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
`, app, primaryRegion, image, internalPort)

	return m.flyBase(token).
		WithNewFile("/fly.toml", config).
		WithExec([]string{"/root/.fly/bin/flyctl", "deploy", "--config", "/fly.toml"}).
		Stdout(ctx)
}

// flyBase returns a container with flyctl installed
func (m *Flyio) flyBase(token *dagger.Secret) *dagger.Container {
	return dag.Alpine().
		Base().
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"sh", "-c", "curl -L https://fly.io/install.sh | sh"}).
		WithSecretVariable("FLY_API_TOKEN", token)
}
