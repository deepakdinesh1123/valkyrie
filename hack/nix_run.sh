#!/usr/bin/env bash

set -e

source ~/.profile
# sed -i "1a #! nix-shell -I nixpkgs=$HOME/24.05" ~/odin/exec.sh
chmod +x ~/odin/exec.sh

while [ ! -f ~/status.txt ]; do sleep 1; done

cd ~/odin
exec ./exec.sh