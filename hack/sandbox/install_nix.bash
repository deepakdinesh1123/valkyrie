#!/usr/bin/env bash
set -euo pipefail

# Default values
NIX_VERSION="2.26.3"
EXTRA_NIX_CONFIG=""
INSTALL_OPTIONS=""
NIX_CHANNEL_REV="b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c" # pragma: allowlist secret

NIXPKGS_OWNER="NixOS"
NIXPKGS_REPO="nixpkgs"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --nix-version)
      NIX_VERSION="$2"
      shift 2
      ;;
    --extra-nix-config)
      EXTRA_NIX_CONFIG="$2"
      shift 2
      ;;
    --install-options)
      INSTALL_OPTIONS="$2"
      shift 2
      ;;
    --nix-channel-rev)
      NIX_CHANNEL_REV="$2"
      shift 2
      ;;
    --help)
      echo "Usage: $0 [OPTIONS]"
      echo "Options:"
      echo "  --nix-version VERSION    Specify Nix version to install (default: 2.24.0)"
      echo "  --extra-nix-config TEXT  Additional Nix configuration"
      echo "  --install-options TEXT   Additional installer options"
      echo "  --nix-channel-rev REV    Specify nixpkgs revision (default: b27ba4e...)"
      echo "  --help                   Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Check if Nix is already installed - wrap in if condition to avoid errors if not found
if command -v nix &> /dev/null; then
  nix_path="$(command -v nix)"
  echo "Aborting: Nix is already installed at ${nix_path}"
  exit 1
fi

echo "Installing Nix version $NIX_VERSION"

# Create a temporary workdir
workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

# Configure Nix
add_config() {
  echo "$1" >> "$workdir/nix.conf"
}
add_config "show-trace = true"
# Set jobs to number of cores
add_config "max-jobs = auto"
# Allow binary caches for root user
add_config "trusted-users = root ${USER:-}"

# Append extra nix configuration if provided
if [[ -n "${EXTRA_NIX_CONFIG}" ]]; then
  add_config "$EXTRA_NIX_CONFIG"
fi
if [[ ! $EXTRA_NIX_CONFIG =~ "experimental-features" ]]; then
  add_config "experimental-features = nix-command flakes"
fi
# Always allow substituting from the cache
if [[ ! $EXTRA_NIX_CONFIG =~ "always-allow-substitutes" ]]; then
  add_config "always-allow-substitutes = true"
fi

# Fix the "error: the group 'nixbld' specified in 'build-users-group' does not exist"
add_config "build-users-group ="

# Nix installer flags - always use no-daemon mode
installer_options=(
  --no-daemon
  --no-channel-add
)

# Add any extra installer options
if [[ -n "${INSTALL_OPTIONS}" ]]; then
  IFS=' ' read -r -a extra_installer_options <<< "$INSTALL_OPTIONS"
  installer_options=("${extra_installer_options[@]}" "${installer_options[@]}")
fi

echo "Installer options: ${installer_options[*]}"

# Download the installer with retries
curl_retries=5
install_url="https://releases.nixos.org/nix/nix-${NIX_VERSION}/install"
echo "Downloading Nix installer from $install_url"

while ! curl -sS -o "$workdir/install" -v --fail -L "$install_url"
do
  sleep 1
  ((curl_retries--))
  if [[ $curl_retries -le 0 ]]; then
    echo "curl retries failed" >&2
    exit 1
  fi
  echo "Retrying download..."
done

# Run the installer
sh "$workdir/install" "${installer_options[@]}"

# Print success message and next steps
echo "Nix $NIX_VERSION has been successfully installed in no-daemon mode!"

export PATH="/home/valnix/.nix-profile/bin:$PATH"

# Update channels
nix-channel --remove nixpkgs
nix-channel --add "https://github.com/$NIXPKGS_OWNER/$NIXPKGS_REPO/archive/$NIX_CHANNEL_REV.tar.gz" nixpkgs
nix-channel --update
