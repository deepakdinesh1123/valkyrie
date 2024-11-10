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
        nix-shell -p "$name" --run "exit"
    else
        nix-shell -p "$language.$name" --run "exit"
    fi
done < "$csv_file"

# tail -n +2 "$csv_file" | while IFS=, read -r name type language
# do
#     # Remove leading/trailing spaces (safe for POSIX shell)
#     name=$(echo "$name" | xargs)
#     type=$(echo "$type" | xargs)
#     language=$(echo "$language" | xargs)

#     if [ "$type" = "system" ]; then
#         echo "Adding $name to store"
#         nix shell -p "$name" --run "exit"
#     else
#         nix-shell -p "$language.$name" --run "exit"
#     fi
# done