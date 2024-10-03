#!/usr/bin/env bash

set -e

NIX_BIN=/nix/store/2nhrwv91g6ycpyxvhmvc0xs8p92wp4bk-nix-2.24.9/bin
nix registry add flake:nixpkgs git+file://~/nixpkgs
echo "setup done" >> ~/status.txt
sleep infinity