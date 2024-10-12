#!/usr/bin/env bash

nix_file="build/package/nix/odin.nix"

sed -i 's/^\s*vendorHash = ".*";/  vendorHash = "";/g' "$nix_file"
out=$(eval "nix build" 2>&1)
vendorHash=$(echo $out | grep -o 'sha256-[^ ]*' | tail -n 1)
echo $vendorHash
old_line='vendorHash = "";'
new_line="vendorHash = \"$vendorHash\";"
sed -i "s|$old_line|$new_line|g" "$nix_file"
