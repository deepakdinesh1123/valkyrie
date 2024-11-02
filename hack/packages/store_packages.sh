#!/bin/bash

# Construct the POSTGRES_URL from environment variables
POSTGRES_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}"

# Check if POSTGRES_URL is set
if [ -z "$POSTGRES_URL" ]; then
  echo "Error: POSTGRES_URL environment variable is not set."
  exit 1
fi

# Fetch packages from the database
packages=$(psql "$POSTGRES_URL" -t -c "SELECT id, name, language;")

# Check if any packages were found
if [ -z "$packages" ]; then
  echo "No packages found in the database."
  exit 0
fi

# Process each package
echo "$packages" | while IFS="|" read -r id name language; do
  name=$(echo "$name" | xargs)
  language=$(echo "$language" | xargs)

  if [ -z "$language" ]; then
    echo "Running nix-shell for $name (type: system)..."
    nix-shell -p "$name" --run "exit"
  else
    echo "Running nix-shell for $language.$name (type: language)..."
    nix-shell -p "$language.$name" --run "exit"
  fi
done

echo "All packages processed successfully."
