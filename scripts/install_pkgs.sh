#!/bin/bash
set -euo pipefail

# Example: Only run in remote environments
if [ "${CLAUDE_CODE_REMOTE:-}" != "true" ]; then
  exit 0
fi

# Install Dagger CLI
mkdir -p "$HOME/.local/bin"
curl -fsSL https://dl.dagger.io/dagger/install.sh | BIN_DIR=$HOME/.local/bin sh

# Install Docker
apt-get update -qq 2>&1 | grep -v "^W:" || true
apt-get install -y docker.io > /dev/null 2>&1

# Start Docker service
service docker start

# Wait for Docker to be ready
sleep 2

# Verify Dagger can communicate with Docker
dagger core version
