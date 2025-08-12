#!/usr/bin/env bash
set -euo pipefail

install_nix() {
  local nix_version="${1:-2.26.3}"
  local extra_nix_config="${2:-}"
  local install_options="${3:-}"
  local nix_channel_rev="${4:-b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c}"
  local nixpkgs_owner="NixOS"
  local nixpkgs_repo="nixpkgs"

  if command -v nix &> /dev/null; then
    local nix_path="$(command -v nix)"
    echo "Aborting: Nix is already installed at ${nix_path}"
    return 0
  fi

  echo "Installing Nix version $nix_version"

  local workdir
  workdir=$(mktemp -d) || {
    echo "ERROR: Failed to create temporary directory using mktemp -d." >&2
    exit 1
  }

  mkdir -p "$workdir"
  touch "$workdir/nix.conf"

  add_config() {
    echo "$1" >> "$workdir/nix.conf"
  }

  add_config "show-trace = true"
  add_config "max-jobs = auto"
  local current_user="${USER:-$(whoami)}"
  if getent passwd "$current_user" > /dev/null; then
      add_config "trusted-users = root $current_user"
  else
      add_config "trusted-users = root $(whoami)"
  fi

  if [[ -n "${extra_nix_config}" ]]; then
    add_config "$extra_nix_config"
  fi

  if ! grep -q "experimental-features" "$workdir/nix.conf"; then
    add_config "experimental-features = nix-command flakes"
  fi

  if ! grep -q "always-allow-substitutes" "$workdir/nix.conf"; then
    add_config "always-allow-substitutes = true"
  fi

  add_config "build-users-group ="

  local installer_options=(
    --no-daemon
    --no-channel-add
    --nix-extra-conf-file "$workdir/nix.conf"
  )

  if [[ -n "${install_options}" ]]; then
      installer_options+=(${install_options})
  fi

  echo "Installer options: ${installer_options[*]}"

  local curl_retries=5
  local install_url="https://releases.nixos.org/nix/nix-${nix_version}/install"
  echo "Downloading Nix installer from $install_url"

  for i in $(seq 1 $curl_retries); do
      if curl -sS -o "$workdir/install" --fail -L "$install_url"; then
          echo "Download successful."
          break
      else
          echo "Download failed (Attempt $i/$curl_retries). Retrying in 1 second..." >&2
          sleep 1
          if [[ $i -eq $curl_retries ]]; then
              echo "ERROR: Failed to download Nix installer after multiple retries" >&2
              exit 1
          fi
      fi
  done

  chmod +x "$workdir/install"

  echo "Running Nix installer..."
  "$workdir/install" "${installer_options[@]}"

  local nix_profile_script="$HOME/.nix-profile/etc/profile.d/nix.sh"
  if [[ -e "$nix_profile_script" ]]; then
    echo "Sourcing Nix profile script: $nix_profile_script"
    . "$nix_profile_script"
  else
      echo "Warning: Nix profile script not found at $nix_profile_script. You may need to source it manually in your shell configuration."
  fi

  export PATH="$HOME/.nix-profile/bin:$PATH"

  if ! command -v nix &> /dev/null; then
    echo "ERROR: Nix installation failed - 'nix' command not found in PATH after install." >&2
    echo "Please check the installer output above for details." >&2
    exit 1
  fi

  echo "Updating Nix channels..."
  if nix-channel --list | grep -q "nixpkgs"; then
    nix-channel --remove nixpkgs
  fi
  nix-channel --add "https://github.com/$nixpkgs_owner/$nixpkgs_repo/archive/$nix_channel_rev.tar.gz" nixpkgs
  nix-channel --update

  echo "Nix $nix_version has been successfully installed in no-daemon mode!"
  echo "Nix channel set to: nixpkgs rev $nix_channel_rev"
  return 0
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  install_nix "${1:-}" "${2:-}" "${3:-}" "${4:-}"
fi
