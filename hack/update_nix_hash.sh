#!/usr/bin/env bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <nix_file> <package_name>" # The package name is the name specified in the nix flake
    exit 1
fi

nix_file="builds/nix/$1.nix"
package_name="$2"

sed -i 's/^\s*vendorHash = ".*";/  vendorHash = "";/g' "$nix_file"
out=$(eval "nix build .#$package_name" 2>&1)
vendorHash=$(echo $out | grep -o 'sha256-[^ ]*' | tail -n 1)
echo $vendorHash
old_line='vendorHash = "";'
new_line="vendorHash = \"$vendorHash\";"
sed -i "s|$old_line|$new_line|g" "$nix_file"
