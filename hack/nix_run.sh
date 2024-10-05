#!/usr/bin/env bash

set -e

source ~/.profile
sed -i "1a #! nix-shell -I nixpkgs=$HOME/24.05" ~/odin/exec.sh
chmod +x ~/odin/exec.sh
cd ~/odin
exec ./exec.sh