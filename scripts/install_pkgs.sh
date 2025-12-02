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
apt-get install -y podman > /dev/null 2>&1

# Start Podman socket service
mkdir -p /run/user/0/podman
podman system service --time=0 unix:///run/user/0/podman/podman.sock &

# Export DOCKER_HOST for Dagger to use Podman
export DOCKER_HOST=unix:///run/user/0/podman/podman.sock
echo "export DOCKER_HOST=unix:///run/user/0/podman/podman.sock" >> "$HOME/.bashrc"
