#!/bin/bash
set -euo pipefail

# Example: Only run in remote environments
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

# Install Dagger CLI
mkdir -p "$HOME/.local/bin"
curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh

# Install Podman
apt-get update -qq 2>&1 | grep -v "^W:" || true

# Install podman with --fix-missing to handle repository issues
apt-get install -y podman --fix-missing > /dev/null 2>&1 || true

# Fix any remaining broken dependencies
apt --fix-broken install -y > /dev/null 2>&1 || true

# Configure Podman for rootful execution (required by Dagger)
# Start rootful Podman socket service
mkdir -p /run/podman
podman system service --time=0 unix:///run/podman/podman.sock > /var/log/podman.log 2>&1 &

# Wait for socket to be ready
sleep 3

# Export DOCKER_HOST for Dagger to use Podman (rootful socket)
export DOCKER_HOST=unix:///run/podman/podman.sock
echo "export DOCKER_HOST=unix:///run/podman/podman.sock" >> "$HOME/.bashrc"

# Verify Podman socket is responding
if ! podman --remote --url unix:///run/podman/podman.sock info > /dev/null 2>&1; then
  echo "Error: Podman socket failed to start. Check /var/log/podman.log for details."
  exit 1
fi

echo "Podman is running successfully"

# Note: Dagger may not be able to pull images from registry.dagger.io
# due to network restrictions in some environments.
# Verify Dagger can communicate with Podman (may fail on first run due to registry access)
echo "Testing Dagger connection to Podman..."
dagger core version || echo "Warning: Dagger cannot pull engine image due to network restrictions"
