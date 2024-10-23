#!/usr/bin/env sh

set -e

mkdir -p ~/.local/state/nix/profiles
ln -s $NIX_CHANNELS_ENVIRONMENT ~/.local/state/nix/profiles/channels-1-link
ln -s ~/.local/state/nix/profiles/channels-1-link ~/.local/state/nix/profiles/channels
ln -s $NIX_USER_ENVIRONMENT ~/.local/state/nix/profiles/profiles-1-link
ln -s ~/.local/state/nix/profiles/profiles-1-link ~/.local/state/nix/profiles/profile

mkdir ~/.nix-defexpr
ln -s ~/.local/state/nix/profiles/channels ~/.nix-defexpr/channels

ln -s $NIX_USER_ENVIRONMENT ~/.nix-profile

echo PATH="$PATH:~/.nix-profile/bin" >> ~/.profile
echo NIX_PATH="/tmp/nixpkgs" >> ~/.profile

# echo export PATH="$PATH:~/.nix-profile/bin" >> ~/.bashrc
# echo export NIX_PATH="/tmp/nixpkgs" >> ~/.bashrc
. ~/.profile

# nix registry add flake:nixpkgs git+file://$HOME/nixpkgs

# nix-channel --add https://nixos.org/channels/nixos-24.05 nixpkgs
# nix-channel --update

echo "setup done" >> ~/status.txt
sleep infinity