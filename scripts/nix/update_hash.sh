#!/usr/bin/env bash
# Updates SRI hashes for flake.nix.

OUT=$(mktemp -d -t nar-hash-XXXXXX)

echo "Downloading Go modules..."
GOPATH="$OUT" go mod download
echo "Calculating SRI hash..."
HASH=$(nardump --sri "$OUT/pkg/mod/cache/download")
sudo rm -rf "$OUT"

packages=("geri" "odin")

for package in "${packages[@]}"; do
    filepath="./build/package/nix/${package}.nix"
    
    if [[ -f "$filepath" ]]; then
        sed -i "s#\(vendorHash = \"\)[^\"]*#\1${HASH}#" "$filepath"
    else
        echo "File $filepath does not exist."
    fi
done