#!/usr/bin/env bash

set -e

mkdir -p ~/.local/state/nix/profiles
ln -s /nix/store/dxq1f3wgkxq5lvv1pn2xa6b15w2f8h27-user-environment ~/.local/state/nix/profiles/channels-1-link
ln -s ~/.local/state/nix/profiles/channels-1-link ~/.local/state/nix/profiles/channels
ln -s /nix/store/ak09zx2rza4x0c4fjn9zyjnr4nck9h1b-user-environment ~/.local/state/nix/profiles/profiles-1-link
ln -s ~/.local/state/nix/profiles/profiles-1-link ~/.local/state/nix/profiles/profile

mkdir ~/.nix-defexpr
ln -s ~.local/state/nix/profiles/channels ~/.nix-defexpr/channels

ln -s /nix/store/ak09zx2rza4x0c4fjn9zyjnr4nck9h1b-user-environment ~/.nix-profile

echo PATH="$PATH:~/.nix-profile/bin" >> ~/.profile
echo export PATH="$PATH:~/.nix-profile/bin" >> ~/.bashrc
source ~/.profile

# not required since we moved to using script
# nix registry add flake:nixpkgs git+file://$HOME/24.05

nix-channel --add https://nixos.org/channels/nixos-24.05 nixpkgs
nix-channel --update

echo "setup done" >> ~/status.txt
sleep infinity