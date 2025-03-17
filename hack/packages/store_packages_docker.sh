#!/bin/sh
# set -e

# Path to your CSV file
csv_file="/tmp/packages.csv"

# Read the file line by line
while IFS=, read -r name type language
do
    # Remove leading/trailing spaces (this is safe for POSIX shell)
    name=$(echo "$name" | xargs)
    type=$(echo "$type" | xargs)
    language=$(echo "$language" | xargs)

    echo $name $type $language
    
    if [ "$type" = "system" ]; then
        cached-nix-shell -p "$name" --run "exit"
    else
        cached-nix-shell -p "$language.$name" --run "exit"
    fi
done < "$csv_file"