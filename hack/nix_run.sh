#!/usr/bin/env bash

set -e

source ~/.bashrc
sed -i "1a #! nix-shell -I nixpkgs=$HOME/24.05" ~/odin/exec.sh
chmod +x ~/odin/exec.sh
exec ~/odin/exec.sh