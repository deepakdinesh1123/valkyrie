#!/bin/bash

if [ ! -d "/home/valnix/.nix-profile" ]; then
    mkdir -p ~/.local/state/nix/profiles
    ln -s $NIX_CHANNELS_ENVIRONMENT ~/.local/state/nix/profiles/channels
    ln -s ~/.local/state/nix/profiles/channels ~/.local/state/nix/profiles/channels-1-link
    ln -s $NIX_USER_ENVIRONMENT ~/.local/state/nix/profiles/profile-1-link
    ln -s ~/.local/state/nix/profiles/profile-1-link ~/.local/state/nix/profiles/profile

    mkdir -p ~/.nix-defexpr
    ln -s ~/.local/state/nix/profiles/channels ~/.nix-defexpr/channels

    ln -s $NIX_USER_ENVIRONMENT ~/.nix-profile

    echo 'export PATH="$PATH:~/.nix-profile/bin"' >> ~/.profile
    echo 'export PATH="$PATH:~/.nix-profile/bin"' >> ~/.bashrc
fi
source ~/.profile
source ~/.bashrc

echo "Printing environment variables"
echo $NIX_CHANNELS_ENVIRONMENT
echo $NIX_USER_ENVIRONMENT

NIX_SECRET_KEY_FILE=/tmp/setup/cache-priv-key.pem nix-serve --listen 0.0.0.0:5000
