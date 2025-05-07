#!/bin/bash

if [ -f "/tmp/setup/.env" ]; then
    echo "Found .env at /tmp/setup/.env"
    cat /tmp/setup/.env
    source "/tmp/setup/.env"
fi

if [ ! -d "/home/valnix/.nix-profile" ]; then
    mkdir -p ~/.local/state/nix/profiles
    ln -s $NIX_CHANNELS_ENVIRONMENT ~/.local/state/nix/profiles/channels
    ln -s ~/.local/state/nix/profiles/channels ~/.local/state/nix/profiles/channels-1-link
    ln -s $NIX_USER_ENVIRONMENT ~/.local/state/nix/profiles/profile-1-link
    ln -s ~/.local/state/nix/profiles/profile-1-link ~/.local/state/nix/profiles/profile

    mkdir -p ~/.nix-defexpr
    ln -s ~/.local/state/nix/profiles/channels ~/.nix-defexpr/channels

    ln -s $NIX_USER_ENVIRONMENT ~/.nix-profile
    echo 'https://github.com/NixOS/nixpkgs/archive/b27ba4eb322d9d2bf2dc9ada9fd59442f50c8d7c.tar.gz nixpkgs' >> ~/.nix-channels
fi

PROFILE_PATH_LINE='export PATH="$PATH:~/.nix-profile/bin"'
BASHRC_PATH_LINE='export PATH="$PATH:~/.nix-profile/bin"'

if ! grep -qF "$PROFILE_PATH_LINE" ~/.profile 2>/dev/null; then
    echo "Adding Nix bin path to ~/.profile"
    echo "$PROFILE_PATH_LINE" >> ~/.profile
else
    echo "Nix bin path already present in ~/.profile"
fi

# Check and add to ~/.bashrc
if ! grep -qF "$BASHRC_PATH_LINE" ~/.bashrc 2>/dev/null; then
    echo "Adding Nix bin path to ~/.bashrc"
    echo "$BASHRC_PATH_LINE" >> ~/.bashrc
else
    echo "Nix bin path already present in ~/.bashrc"
fi

source ~/.profile
source ~/.bashrc

echo "PATH is $PATH"

NIX_SECRET_KEY_FILE=/home/valnix/.config/cache-priv-key.pem nix-serve --listen 0.0.0.0:5000
