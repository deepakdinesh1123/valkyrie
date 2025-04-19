#!/usr/bin/env bash
set -euo pipefail

install_nix() {
  # Set default values first
  local nix_version="${1:-2.26.3}"
  local extra_nix_config="${2:-}"
  local install_options="${3:-}"
  local nix_channel_rev="${4:-b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c}"
  local nixpkgs_owner="NixOS"
  local nixpkgs_repo="nixpkgs"

  # Check if Nix is already installed
  if command -v nix &> /dev/null; then
    local nix_path="$(command -v nix)"
    echo "Aborting: Nix is already installed at ${nix_path}"
    exit 0
  fi

  echo "Installing Nix version $nix_version"

  # Create a temporary workdir
  local workdir=$(mktemp -d)
  trap 'rm -rf "$workdir"' EXIT

  # Configure Nix
  mkdir -p "$workdir"
  touch "$workdir/nix.conf"
  
  add_config() {
    echo "$1" >> "$workdir/nix.conf"
  }
  
  add_config "show-trace = true"
  # Set jobs to number of cores
  add_config "max-jobs = auto"
  # Allow binary caches for root user
  add_config "trusted-users = root ${USER:-$(whoami)}"

  # Append extra nix configuration if provided
  if [[ -n "${extra_nix_config}" ]]; then
    add_config "$extra_nix_config"
  fi
  
  # Add experimental features if not already specified
  if [[ ! "${extra_nix_config}" =~ "experimental-features" ]]; then
    add_config "experimental-features = nix-command flakes"
  fi
  
  # Always allow substituting from the cache if not already specified
  if [[ ! "${extra_nix_config}" =~ "always-allow-substitutes" ]]; then
    add_config "always-allow-substitutes = true"
  fi

  # Fix the "error: the group 'nixbld' specified in 'build-users-group' does not exist"
  add_config "build-users-group ="

  # Nix installer flags - always use no-daemon mode
  local installer_options=(
    --no-daemon
    --no-channel-add
    --nix-extra-conf-file "$workdir/nix.conf"
  )

  # Add any extra installer options
  if [[ -n "${install_options}" ]]; then
    IFS=' ' read -r -a extra_installer_options <<< "$install_options"
    installer_options+=("${extra_installer_options[@]}")
  fi

  echo "Installer options: ${installer_options[*]}"

  # Download the installer with retries
  local curl_retries=5
  local install_url="https://releases.nixos.org/nix/nix-${nix_version}/install"
  echo "Downloading Nix installer from $install_url"

  while ! curl -sS -o "$workdir/install" --fail -L "$install_url"
  do
    sleep 1
    ((curl_retries--))
    if [[ $curl_retries -le 0 ]]; then
      echo "ERROR: Failed to download Nix installer after multiple retries" >&2
      exit 1
    fi
    echo "Retrying download..."
  done

  chmod +x "$workdir/install"

  # Run the installer
  echo "Running Nix installer..."
  "$workdir/install" "${installer_options[@]}"

  # Source the nix environment if it exists
  if [[ -e "$HOME/.nix-profile/etc/profile.d/nix.sh" ]]; then
    source "$HOME/.nix-profile/etc/profile.d/nix.sh"
  fi

  # Ensure nix is in PATH
  export PATH="$HOME/.nix-profile/bin:$PATH"

  # Verify nix was installed correctly
  if ! command -v nix &> /dev/null; then
    echo "ERROR: Nix installation failed - 'nix' command not found in PATH" >&2
    exit 1
  fi

  # Update channels
  echo "Updating Nix channels..."
  nix-channel --remove nixpkgs
  nix-channel --add "https://github.com/$nixpkgs_owner/$nixpkgs_repo/archive/$nix_channel_rev.tar.gz" nixpkgs
  nix-channel --update

  # Print success message
  echo "Nix $nix_version has been successfully installed in no-daemon mode!"
  echo "Nix channel set to: nixpkgs rev $nix_channel_rev"
  echo "To use Nix in this session, run: source $HOME/.nix-profile/etc/profile.d/nix.sh"
  
  exit 0
}

# If executed as a script (not sourced), call the function with arguments
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  install_nix "${1:-}" "${2:-}" "${3:-}" "${4:-}"
fi
