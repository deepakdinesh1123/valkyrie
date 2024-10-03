#!/usr/bin/env bash

source ~/.bashrc
NIX_BIN=/nix/store/2nhrwv91g6ycpyxvhmvc0xs8p92wp4bk-nix-2.24.9/bin
FILE=~/status.txt

while [ ! -f "$FILE" ]; do
    sleep 1
done
cd ~/odin
$NIX_BIN/nix run