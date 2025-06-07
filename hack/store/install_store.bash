#!/usr/bin/env bash
set -euo pipefail

ENVIRONMENT="${1:-container}"

# Check if nix-serve is already installed
if command -v nix-serve &> /dev/null; then
    NIX_SERVE_PATH="$(command -v nix-serve)"
    echo "Aborting: Nix-Serve is already installed at ${NIX_SERVE_PATH}"
    exit 0
fi

# Source the install_nix script
if [ ! -f ./install_nix.bash ]; then
    echo "Error: install_nix.bash not found!"
    exit 1
fi
source ./install_nix.bash

# Verify the install_nix function exists
if ! declare -f install_nix &> /dev/null; then
    echo "Error: install_nix function not found in install_nix.bash!"
    exit 1
fi

# Install Nix
install_nix

# Check if valnix user exists
if ! id "valnix" &>/dev/null; then
    echo "Error: User 'valnix' does not exist. Cannot chown directories."
    exit 1
fi

# Create directories and set permissions
echo "Creating directories and setting permissions..."
sudo chown -R valnix:valnix /nix
sudo chown -R valnix:valnix /tmp/setup

# Install nix-serve-ng
echo "Installing nix-serve-ng..."
if ! nix-env -iA nixpkgs.nix-serve-ng; then
    echo "Error: Failed to install nix-serve-ng."
    exit 1
fi

# Generate binary cache keys if not in k8s environment
if [[ "${ENVIRONMENT}" != "k8s" ]]; then
    echo "Generating binary cache keys..."
    if nix-store --generate-binary-cache-key valkyrie-store /tmp/setup/cache-priv-key.pem /tmp/setup/cache-pub-key.pem; then
        echo "Generated binary cache keys successfully."
        chmod 600 /tmp/setup/cache-priv-key.pem
    else
        echo "Error: Failed to generate binary cache keys."
        exit 1
    fi
else
    echo "Skipping binary cache key generation for Kubernetes environment."
fi

# Create .env file
echo "Adding channels and user environment to .env"
touch /tmp/setup/.env
sudo chown valnix:valnix /tmp/setup/.env

# Add Nix channels to .env
NIX_CHANNELS_PATH="$HOME/.local/state/nix/profiles/channels"
if [ -L "$NIX_CHANNELS_PATH" ] || [ -d "$NIX_CHANNELS_PATH" ]; then
    REAL_PATH=$(realpath "$NIX_CHANNELS_PATH")
    bash -c "echo \"NIX_CHANNELS_ENVIRONMENT=$REAL_PATH\" >> /tmp/setup/.env"
else
    echo "Warning: Nix channels path not found at $NIX_CHANNELS_PATH"
fi

# Add Nix user environment to .env
NIX_PROFILE_PATH="$HOME/.nix-profile"
if [ -L "$NIX_PROFILE_PATH" ] || [ -d "$NIX_PROFILE_PATH" ]; then
    REAL_PATH=$(realpath "$NIX_PROFILE_PATH")
    bash -c "echo \"NIX_USER_ENVIRONMENT=$REAL_PATH\" >> /tmp/setup/.env"
else
    echo "Warning: Nix user environment path not found at $NIX_PROFILE_PATH"
fi
