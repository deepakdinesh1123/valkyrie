#!/bin/bash

# Check if the required arguments are passed
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <language> <schema_folder> <target_folder>"
    exit 1
fi

LANGUAGE="$1"
SCHEMA_FOLDER="$2"
TARGET_FOLDER="$3"

for file in "$SCHEMA_FOLDER"/*.json; do
    # Check if there are JSON files in the directory
    if [ ! -e "$file" ]; then
        echo "No JSON files found in '$SCHEMA_FOLDER'."
        break
    fi

    FILENAME=$(basename "$file" .json)

    case $LANGUAGE in
        go)
            quicktype -o agent/schema/$FILENAME.go --src-lang schema --package schema --lang go $file 
            ;;
        *)
            echo "Specified language $LANGUAGE is not supported"
            exit 1
            ;;
    esac

    
done

