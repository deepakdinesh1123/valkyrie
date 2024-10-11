#!/bin/bash

export PGPASSWORD="thorkell"

for path in $(psql -U thors -d valkyrie -h localhost -t -c "SELECT path FROM packages WHERE arch IS NULL LIMIT 10;"); do
  if nix-store --realise "$path"; then
    arch=$(uname -m)-$(uname -s | tr '[:upper:]' '[:lower:]')
    
    psql -U thors -d valkyrie -h localhost -c "UPDATE packages SET arch = '$arch' WHERE path = '$path';"
  else
    echo "Failed to install $path - possibly different architecture."
  fi
done

unset PGPASSWORD
