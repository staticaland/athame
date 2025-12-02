#!/bin/bash
set -euo pipefail

# Example: Only run in remote environments
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

# Install Dagger CLI
mkdir -p "$HOME/.local/bin"
curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh

# Install Docker with specific runc version to avoid 404 errors
apt-get update -qq 2>&1 | grep -v "^W:" || true

# Install older runc version first to avoid 404 errors on newer version
apt-get install -y runc=1.1.12-0ubuntu3 > /dev/null 2>&1 || true

# Install docker.io with --fix-missing to handle repository issues
apt-get install -y docker.io --fix-missing > /dev/null 2>&1 || true

# Fix any remaining broken dependencies
apt --fix-broken install -y > /dev/null 2>&1 || true

# Start Docker daemon with flags to work around networking restrictions
# --iptables=false: Skip iptables setup (required in restricted environments)
# --ip-masq=false: Disable IP masquerading
# --bridge=none: Don't create default bridge network
dockerd --iptables=false --ip-masq=false --bridge=none > /var/log/dockerd.log 2>&1 &

# Wait for Docker to be ready
sleep 5

# Verify Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Error: Docker failed to start. Check /var/log/dockerd.log for details."
  exit 1
fi

echo "Docker is running successfully"

# Note: Dagger may not be able to pull images from registry.dagger.io
# due to network restrictions in some environments.
# Verify Dagger can communicate with Docker (may fail on first run due to registry access)
echo "Testing Dagger connection to Docker..."
dagger core version || echo "Warning: Dagger cannot pull engine image due to network restrictions"
