#!/usr/bin/env bash
set -euo pipefail

ENVIRONMENT="${1:-container}"

# Check if Nix-Serve is already installed
if command -v nix-serve &> /dev/null; then
    local nix_path="$(command -v nix-serve)"
    echo "Aborting: Nix-Serve is already installed at ${nix_path}"
    exit 0
fi

source ./install_nix.bash

install_nix

sudo chown -R valnix:valnix /nix
sudo chown -R valnix:valnix /tmp/setup

nix-env -iA nixpkgs.nix-serve-ng
if [[ "${ENVIRONMENT}" != "k8s" ]]; then
    echo "Generating binary cache keys..."
    nix-store --generate-binary-cache-key valkyrie-store /tmp/setup/cache-priv-key.pem /tmp/setup/cache-pub-key.pem && \
        chown valnix /tmp/setup/cache-priv-key.pem && \
        chmod 600 /tmp/setup/cache-priv-key.pem
else
    echo "Skipping binary cache key generation for Kubernetes environment."
fi

echo "NIX_CHANNELS_ENVIRONMENT=$(realpath ~/.local/state/nix/profiles/channels)" >> /tmp/setup/.env
echo "NIX_USER_ENVIRONMENT=$(realpath ~/.nix-profile/)" >> /tmp/setup/.env
